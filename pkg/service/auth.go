package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mi-gongan/commention_backend/pkg/constant"
	"github.com/mi-gongan/commention_backend/pkg/model"
)

func CreateToken(user model.UserForJWT) (string, error) {
	expirationTime := time.Now().Add(24 * 60 * time.Minute)

	claims := &model.Claims{
		UserForJWT: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(constant.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateTokenWithRefresh(user model.UserForJWT) (string, string, error) {
	tokenString, err := CreateToken(user)

	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenString, err := refreshToken.SignedString(constant.SecretKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func VerifyToken(tokenString string) (*model.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return constant.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func VerifyRefreshToken(refreshTokenString string) (*model.Claims, error) {

	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return constant.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := refreshToken.Claims.(*model.Claims); ok && refreshToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
