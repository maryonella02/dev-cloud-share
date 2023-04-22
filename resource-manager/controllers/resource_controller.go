package controllers

import (
	"dev-cloud-share/resource-manager/models"
	"dev-cloud-share/resource-manager/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ResourceController struct {
	resourceService *services.ResourceService
}

func NewResourceController(resourceService *services.ResourceService) *ResourceController {
	return &ResourceController{resourceService}
}

func (rc *ResourceController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/resources", rc.registerResource).Methods("POST")
	router.HandleFunc("/resources", rc.getResources).Methods("GET")
	router.HandleFunc("/resources/{resource_id}", rc.updateResource).Methods("PUT")
	router.HandleFunc("/resources/{resource_id}", rc.deleteResource).Methods("DELETE")
	router.HandleFunc("/allocations", rc.allocateResource).Methods("POST")
	router.HandleFunc("/allocations/{allocation_id}", rc.releaseResource).Methods("DELETE")
}

func (rc *ResourceController) registerResource(w http.ResponseWriter, r *http.Request) {
	var resource models.Resource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rc.resourceService.RegisterResource(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resource)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (rc *ResourceController) getResources(w http.ResponseWriter, r *http.Request) {
	resources, err := rc.resourceService.GetResources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resources)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (rc *ResourceController) updateResource(w http.ResponseWriter, r *http.Request) {
	var updatedResource models.Resource
	err := json.NewDecoder(r.Body).Decode(&updatedResource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resourceID := mux.Vars(r)["resource_id"]
	err = rc.resourceService.UpdateResource(resourceID, &updatedResource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *ResourceController) deleteResource(w http.ResponseWriter, r *http.Request) {
	resourceID := mux.Vars(r)["resource_id"]
	err := rc.resourceService.DeleteResource(resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *ResourceController) allocateResource(w http.ResponseWriter, r *http.Request) {
	var allocationInfo struct {
		BorrowerID string `json:"borrower_id"`
		ResourceID string `json:"resource_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&allocationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rc.resourceService.AllocateResource(allocationInfo.BorrowerID, allocationInfo.ResourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *ResourceController) releaseResource(w http.ResponseWriter, r *http.Request) {
	allocationID := mux.Vars(r)["allocation_id"]
	err := rc.resourceService.ReleaseResource(allocationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
