package model

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Login  string `json:"login"`
	Role   int    `json:"role"`
	UserID string `json:"userID"`
	jwt.StandardClaims
}
