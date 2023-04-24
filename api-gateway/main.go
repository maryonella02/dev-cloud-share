package main

import (
	"context"
	"dev-cloud-share/resource-manager/controllers"
	"dev-cloud-share/resource-manager/services"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	resourceManagerBaseURL = "http://localhost:8080"
)

func main() {
	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("resource-manager")
	resourceService := services.NewResourceService(db)
	resourceController := controllers.NewResourceController(resourceService)

	// Setup the API Gateway routes
	r := mux.NewRouter()

	// Resource Manager routes
	r.HandleFunc("/resources", resourceController.RegisterResource).Methods("POST")
	r.HandleFunc("/resources", resourceController.GetResources).Methods("GET")
	r.HandleFunc("/resources/{resource_id}", resourceController.UpdateResource).Methods("PUT")
	r.HandleFunc("/resources/{resource_id}", resourceController.DeleteResource).Methods("DELETE")
	r.HandleFunc("/allocations", resourceController.AllocateResource).Methods("POST")
	r.HandleFunc("/allocations/{allocation_id}", resourceController.ReleaseResource).Methods("DELETE")

	// TODO: Add more routes for other microservices

	// Start the API Gateway server
	fmt.Println("API Gateway listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
