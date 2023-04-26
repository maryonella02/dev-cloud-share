package controllers

import (
	"dev-cloud-share/containerization-engine/models"
	"dev-cloud-share/containerization-engine/services"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ContainerController struct {
	containerService *services.ContainerService
}

func NewContainerController(containerService *services.ContainerService) *ContainerController {
	return &ContainerController{containerService}
}

func (c *ContainerController) CreateContainer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var containerConfig models.ContainerConfig
	if err := json.Unmarshal(body, &containerConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	containerID, err := c.containerService.CreateContainer(containerConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{
		ID: containerID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ContainerController) StartContainer(w http.ResponseWriter, r *http.Request) {
	containerID := mux.Vars(r)["id"]

	if err := c.containerService.StartContainer(containerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ContainerController) StopContainer(w http.ResponseWriter, r *http.Request) {
	containerID := mux.Vars(r)["id"]

	if err := c.containerService.StopContainer(containerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ContainerController) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	containerID := mux.Vars(r)["id"]

	if err := c.containerService.RemoveContainer(containerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ContainerController) GetContainerStatus(w http.ResponseWriter, r *http.Request) {
	containerID := mux.Vars(r)["id"]

	status, err := c.containerService.GetContainerStatus(containerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ContainerController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/containers", c.CreateContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/start", c.StartContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/stop", c.StopContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/remove", c.RemoveContainer).Methods("POST")
	router.HandleFunc("/containers/{id}/status", c.GetContainerStatus).Methods("GET")
}
