package controllers

import (
	"auth/models"
	"auth/services"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var jwtKey = []byte("your_secret_key")

type Controller struct {
	authService *services.Service
}

func NewAuthController(authService *services.Service) *Controller {
	return &Controller{authService}
}

func (ac *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", ac.RegisterUser).Methods("POST")
	router.HandleFunc("/login", ac.LoginUser).Methods("POST")
}

func (ac *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registeredUser, err := ac.authService.RegisterUser(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear the password before sending the response
	registeredUser.Password = ""

	_, err = json.Marshal(registeredUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (ac *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ac.authService.LoginUser(r.Context(), loginInfo.Username, loginInfo.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Create the JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = user.Username
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
