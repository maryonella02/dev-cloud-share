package services

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
)

var JWTKey = "secret"

type AuthService struct {
	ServiceURL string
}

func NewAuthService(serviceURL string) *AuthService {
	return &AuthService{
		ServiceURL: serviceURL,
	}
}

func (as *AuthService) ProxyRequestWithToken(w http.ResponseWriter, r *http.Request, endpoint, method string) (string, error) {
	target := as.ServiceURL + endpoint
	req, err := http.NewRequest(method, target, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	// Copy headers from original request to the new request
	for k, v := range r.Header {
		req.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Return original response
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return "", fmt.Errorf("please try again")
	}

	var response struct {
		Email string `json:"email"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Email, nil
}

func (as *AuthService) CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	tokenString, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (as *AuthService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTKey), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}
