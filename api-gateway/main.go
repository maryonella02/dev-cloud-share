package main

import (
	"api-gateway/controllers"
	"api-gateway/services"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	resourceManagerBaseURL        = "http://resource-manager:8080"
	containerizationEngineBaseURL = "http://containerization-engine:8082"
	authServiceBaseURL            = "https://auth:8443"
)

func main() {
	router := mux.NewRouter()

	authService := services.NewAuthService(authServiceBaseURL)
	authController := controllers.NewAuthController(authService)
	authController.RegisterRoutes(router.PathPrefix("/api/v1").Subrouter())

	resourceService := services.NewResourceService(resourceManagerBaseURL)
	resourceController := controllers.NewResourceController(resourceService)
	resourceSubRouter := router.PathPrefix("/api/v1").Subrouter()
	resourceSubRouter.Use(authController.JwtAuthMiddleware) // Apply the middleware to resource routes

	resourceController.RegisterRoutes(resourceSubRouter)

	containerizationEngineService := services.NewContainerizationEngineService(containerizationEngineBaseURL)
	containerizationEngineController := controllers.NewContainerizationEngineController(containerizationEngineService)
	containerizationEngineController.RegisterRoutes(router.PathPrefix("/api/v1").Subrouter())

	// Start the server using TLS
	err := http.ListenAndServeTLS(":8440", "./certs/cert.pem", "./certs/key.pem", router)

	if err != nil {
		fmt.Println(err)
		log.Fatal("Server error:", err)
	}
}
