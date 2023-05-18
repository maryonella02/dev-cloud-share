package controllers

import (
	"api-gateway/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService}
}

func (ac *AuthController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", ac.Login).Methods("POST")
	router.HandleFunc("/register", ac.Register).Methods("POST")
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	email, err := ac.authService.ProxyRequestWithToken(w, r, "/login", "POST", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := ac.authService.CreateToken(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var registrationData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&registrationData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	registerData := map[string]any{
		"email":    registrationData.Email,
		"password": registrationData.Password,
	}
	email, err := ac.authService.ProxyRequestWithToken(w, r, "/register", "POST", registerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := ac.authService.CreateToken(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (ac *AuthController) JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore login and register routes
		if strings.Contains(r.URL.Path, "/login") || strings.Contains(r.URL.Path, "/register") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Invalid auth token", http.StatusUnauthorized)
			return
		}

		valid, err := ac.authService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !valid {
			http.Error(w, "Invalid auth token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
