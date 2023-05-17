package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}
