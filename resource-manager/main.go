package main

import (
	"context"
	"dev-cloud-share/resource-manager/controllers"
	"dev-cloud-share/resource-manager/services"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("resource_manager")

	// Create Resource Service and Resource Controller instances
	resourceService := services.NewResourceService(db)
	resourceController := controllers.NewResourceController(resourceService)

	// Create indexes for better performance
	createIndexes(db)

	// Set up the HTTP server and routes
	router := mux.NewRouter()
	resourceController.RegisterRoutes(router)

	log.Println("Resource Manager is running on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func createIndexes(db *mongo.Database) {
	// Index for the 'lender_id' field in the 'resources' collection
	resourcesCollection := db.Collection("resources")
	lenderIDIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"lender_id": 1, // 1 for ascending order
		},
		Options: options.Index().SetName("lender_id_index"),
	}
	_, err := resourcesCollection.Indexes().CreateOne(context.Background(), lenderIDIndexModel)
	if err != nil {
		log.Fatalf("Failed to create index for lender_id: %v", err)
	}

	// Index for the 'resources' field in the 'borrowers' collection
	borrowersCollection := db.Collection("borrowers")
	resourcesIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"resources": 1, // 1 for ascending order
		},
		Options: options.Index().SetName("resources_index"),
	}
	_, err = borrowersCollection.Indexes().CreateOne(context.Background(), resourcesIndexModel)
	if err != nil {
		log.Fatalf("Failed to create index for resources: %v", err)
	}
}
