package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"resource-manager/models"
	"resource-manager/services"
)

type BorrowerController struct {
	resourceService *services.ResourceService
}

func NewBorrowerController(rs *services.ResourceService) *BorrowerController {
	return &BorrowerController{resourceService: rs}
}

func (bc *BorrowerController) CreateBorrower(w http.ResponseWriter, r *http.Request) {
	var borrower models.Borrower
	err := json.NewDecoder(r.Body).Decode(&borrower)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newBorrower, err := bc.resourceService.CreateBorrower(borrower)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newBorrower)
	if err != nil {
		fmt.Println(err)
		return
	}
}
