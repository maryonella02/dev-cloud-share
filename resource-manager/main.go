package main

import (
	"context"
	"dev-cloud-share/resource-manager/controllers"
	"dev-cloud-share/resource-manager/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("resource_manager")

	// Create Resource Service and Resource Controller instances
	resourceService := services.NewResourceService(db)
	resourceController := controllers.NewResourceController(resourceService)

	// Set up the HTTP server and routes
	router := mux.NewRouter()
	resourceController.RegisterRoutes(router)

	log.Println("Resource Manager is running on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
