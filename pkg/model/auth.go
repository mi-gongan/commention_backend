package model

import "github.com/dgrijalva/jwt-go"

type UserForJWT struct {
	Email string `json:"email"`
}

type Claims struct {
	UserForJWT `json:"user"`
	jwt.StandardClaims
}
