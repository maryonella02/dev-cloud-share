package controllers

import (
	"dev-cloud-share/api-gateway/services"
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
	router.HandleFunc("/resources", rc.registerResource).Methods("POST")
	router.HandleFunc("/resources", rc.getResources).Methods("GET")
	router.HandleFunc("/resources/{resource_id}", rc.updateResource).Methods("PUT")
	router.HandleFunc("/resources/{resource_id}", rc.deleteResource).Methods("DELETE")
	router.HandleFunc("/allocations", rc.allocateResource).Methods("POST")
	router.HandleFunc("/allocations/{allocation_id}", rc.releaseResource).Methods("DELETE")
	router.HandleFunc("/borrowers", rc.createBorrower).Methods("POST")
	router.HandleFunc("/lenders", rc.createLender).Methods("POST")
}

func (rc *ResourceController) registerResource(w http.ResponseWriter, r *http.Request) {
	endpoint := "/resources"
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (rc *ResourceController) getResources(w http.ResponseWriter, r *http.Request) {
	endpoint := "/resources"
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "GET")
	if err != nil {
		return
	}
}

func (rc *ResourceController) updateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]
	endpoint := fmt.Sprintf("/resources/%s", resourceID)
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "PUT")
	if err != nil {
		return
	}
}

func (rc *ResourceController) deleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]
	endpoint := fmt.Sprintf("/resources/%s", resourceID)
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "DELETE")
	if err != nil {
		return
	}
}

func (rc *ResourceController) allocateResource(w http.ResponseWriter, r *http.Request) {
	endpoint := "/allocations"
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}

func (rc *ResourceController) releaseResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	allocationID := vars["allocation_id"]
	endpoint := fmt.Sprintf("/allocations/%s", allocationID)
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "DELETE")
	if err != nil {
		return
	}
}

func (rc *ResourceController) createBorrower(w http.ResponseWriter, r *http.Request) {
	endpoint := "/borrowers"
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}

func (rc *ResourceController) createLender(w http.ResponseWriter, r *http.Request) {
	endpoint := "/lenders"
	err := rc.resourceService.ProxyRequest(w, r, endpoint, "POST")
	if err != nil {
		return
	}
}
