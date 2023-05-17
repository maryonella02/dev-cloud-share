package main

import (
	"containerization-engine/controllers"
	"containerization-engine/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	containerService, err := services.NewContainerService()
	if err != nil {
		log.Fatalf("Failed to initialize container service: %v", err)
	}

	containerController := controllers.NewContainerController(containerService)
	containerController.RegisterRoutes(router)

	log.Println("Starting containerization engine on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
