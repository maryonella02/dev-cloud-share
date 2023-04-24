package main

import (
	"dev-cloud-share/api-gateway/controllers"
	"dev-cloud-share/api-gateway/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	resourceManagerBaseURL = "http://localhost:8080"
)

func main() {
	resourceService := services.NewResourceService(resourceManagerBaseURL)
	resourceController := controllers.NewResourceController(resourceService)

	router := mux.NewRouter()

	resourceController.RegisterRoutes(router.PathPrefix("/api/v1").Subrouter())

	err := http.ListenAndServe("localhost:8081", router)
	if err != nil {
		fmt.Println(err)
		return
	}

}
