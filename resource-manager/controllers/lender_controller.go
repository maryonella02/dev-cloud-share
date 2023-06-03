package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"resource-manager/models"
	"resource-manager/services"
)

type LenderController struct {
	resourceService *services.ResourceService
}

func NewLenderController(rs *services.ResourceService) *LenderController {
	return &LenderController{resourceService: rs}
}

func (bc *LenderController) CreateLender(w http.ResponseWriter, r *http.Request) {
	var lender models.Lender
	err := json.NewDecoder(r.Body).Decode(&lender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newLender, err := bc.resourceService.CreateLender(lender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newLender)
	if err != nil {
		fmt.Println(err)
		return
	}
}
