package model

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}
