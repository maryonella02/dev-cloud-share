package controllers

import (
	"dev-cloud-share/api-gateway/services"
	"github.com/gorilla/mux"
	"net/http"
)

type ContainerizationEngineController struct {
	containerizationEngineService *services.ContainerizationEngineService
}

func NewContainerizationEngineController(containerizationEngineService *services.ContainerizationEngineService) *ContainerizationEngineController {
	return &ContainerizationEngineController{containerizationEngineService}
}

func (cec *ContainerizationEngineController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/containers", cec.createContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/start", cec.startContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/stop", cec.stopContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/remove", cec.removeContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/status", cec.getContainerStatus).Methods("GET")
}

func (cec *ContainerizationEngineController) createContainer(w http.ResponseWriter, r *http.Request) {
	endpoint := "/containers"
	err := cec.containerizationEngineService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}

func (cec *ContainerizationEngineController) startContainer(w http.ResponseWriter, r *http.Request) {
	endpoint := "/containers/" + mux.Vars(r)["id"] + "/start"
	err := cec.containerizationEngineService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}

func (cec *ContainerizationEngineController) stopContainer(w http.ResponseWriter, r *http.Request) {
	endpoint := "/containers/" + mux.Vars(r)["id"] + "/stop"
	err := cec.containerizationEngineService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}

func (cec *ContainerizationEngineController) removeContainer(w http.ResponseWriter, r *http.Request) {
	endpoint := "/containers/" + mux.Vars(r)["id"] + "/remove"
	err := cec.containerizationEngineService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}

func (cec *ContainerizationEngineController) getContainerStatus(w http.ResponseWriter, r *http.Request) {
	endpoint := "/containers/" + mux.Vars(r)["id"] + "/status"
	err := cec.containerizationEngineService.ProxyRequest(w, r, endpoint, "GET")
	if err != nil {
		return
	}
}
