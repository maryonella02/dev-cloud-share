package main

import (
	"context"
	"dev-cloud-share/resource-manager/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Define the MongoDB configuration
var mongoURI = "mongodb://localhost:27017"
var dbName = "resource-manager"
var resourceCollection = "resources"
var requestCollection = "requests"

// Define the API endpoints
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/resources", createResource).Methods("POST")
	router.HandleFunc("/resources", getAllResources).Methods("GET")
	router.HandleFunc("/resources/{id}", getResource).Methods("GET")
	router.HandleFunc("/resources/{id}", updateResource).Methods("PUT")
	router.HandleFunc("/requests", createRequest).Methods("POST")
	router.HandleFunc("/requests", getAllRequests).Methods("GET")
	router.HandleFunc("/requests/{id}", getRequest).Methods("GET")
	router.HandleFunc("/requests/{id}", updateRequest).Methods("PUT")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Define the API endpoint handlers
func createResource(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a Resource object
	var resource models.Resource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the resource to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	collection := client.Database(dbName).Collection(resourceCollection)
	result, err := collection.InsertOne(ctx, resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the resource ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)
}
