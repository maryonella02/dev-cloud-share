package controllers

import (
	"dev-cloud-share/resource-manager/models"
	"dev-cloud-share/resource-manager/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ResourceController struct {
	resourceService *services.ResourceService
}

func NewResourceController(resourceService *services.ResourceService) *ResourceController {
	return &ResourceController{resourceService}
}

func (rc *ResourceController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/resources", rc.RegisterResource).Methods("POST")
	router.HandleFunc("/resources", rc.GetResources).Methods("GET")
	router.HandleFunc("/resources/{resource_id}", rc.UpdateResource).Methods("PUT")
	router.HandleFunc("/resources/{resource_id}", rc.DeleteResource).Methods("DELETE")
	router.HandleFunc("/allocations", rc.AllocateResource).Methods("POST")
	router.HandleFunc("/allocations/{resource_id}", rc.ReleaseResource).Methods("DELETE")
}

func (rc *ResourceController) RegisterResource(w http.ResponseWriter, r *http.Request) {
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

func (rc *ResourceController) GetResources(w http.ResponseWriter, r *http.Request) {
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

func (rc *ResourceController) UpdateResource(w http.ResponseWriter, r *http.Request) {
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

func (rc *ResourceController) DeleteResource(w http.ResponseWriter, r *http.Request) {
	resourceID := mux.Vars(r)["resource_id"]
	err := rc.resourceService.DeleteResource(resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *ResourceController) AllocateResource(w http.ResponseWriter, r *http.Request) {
	var allocationInfo struct {
		BorrowerID      string                   `json:"borrower_id"`
		ResourceRequest services.ResourceRequest `json:"resource_request"`
	}
	// example of payload
	//{
	//	"borrower_id": "60a7f8370abf2f3b903bdbb0",
	//		"resource_request": {
	//		"resource_type": "CPU",
	//			"min_cpu_cores": 4,
	//			"min_memory_mb": 1024,
	//			"min_storage_gb": 50
	//	}
	//}

	err := json.NewDecoder(r.Body).Decode(&allocationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := rc.resourceService.AllocateResource(allocationInfo.BorrowerID, allocationInfo.ResourceRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (rc *ResourceController) ReleaseResource(w http.ResponseWriter, r *http.Request) {
	resourceID := mux.Vars(r)["resource_id"]

	releasedResource, err := rc.resourceService.ReleaseResource(resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Resource *models.Resource `json:"resource"`
	}{
		Resource: releasedResource,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
		return
	}
}
