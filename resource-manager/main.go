package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"resource-manager/controllers"
	"resource-manager/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a MongoDB connection string
var connectionString = fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))

func main() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("resource_manager")

	// Create Resource Service and Resource Controller instances
	resourceService := services.NewResourceService(db)
	resourceController := controllers.NewResourceController(resourceService)

	borrowerController := controllers.NewBorrowerController(resourceService)
	lenderController := controllers.NewLenderController(resourceService)

	// Create indexes for better performance
	createIndexes(db)

	// Set up the HTTP server and routes
	router := mux.NewRouter()
	resourceController.RegisterRoutes(router)
	router.HandleFunc("/borrowers", borrowerController.CreateBorrower).Methods(http.MethodPost)
	router.HandleFunc("/lenders", lenderController.CreateLender).Methods(http.MethodPost)

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
