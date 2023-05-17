package main

import (
	"api-gateway/controllers"
	"api-gateway/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	resourceManagerBaseURL        = "http://localhost:8080"
	containerizationEngineBaseURL = "http://localhost:8082"
)

func main() {
	router := mux.NewRouter()

	resourceService := services.NewResourceService(resourceManagerBaseURL)
	resourceController := controllers.NewResourceController(resourceService)
	resourceController.RegisterRoutes(router.PathPrefix("/api/v1").Subrouter())

	containerizationEngineService := services.NewContainerizationEngineService(containerizationEngineBaseURL)
	containerizationEngineController := controllers.NewContainerizationEngineController(containerizationEngineService)
	containerizationEngineController.RegisterRoutes(router.PathPrefix("/api/v1").Subrouter())

	err := http.ListenAndServe("localhost:8081", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
