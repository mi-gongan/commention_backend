package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mi-gongan/commention_backend/pkg/model"
)

var SecretKey []byte

func CreateToken(user model.User) (string, error) {
	// 토큰 만료 시간 설정
	expirationTime := time.Now().Add(5 * time.Minute)

	// 토큰 생성
	claims := &model.Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateTokenWithRefresh(user model.User) (string, string, error) {
	tokenString, err := CreateToken(user)

	if err != nil {
		return "", "", err
	}

	// 리프레시 토큰 생성
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenString, err := refreshToken.SignedString(SecretKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// VerifyToken 함수는 주어진 JWT 토큰을 검증합니다.
func VerifyToken(tokenString string) (*model.Claims, error) {
	// JWT 토큰을 파싱합니다.
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 토큰이 유효하면 클레임을 반환합니다.
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func VerifyRefreshToken(refreshTokenString string) (*model.Claims, error) {
	// 리프레시 토큰을 파싱합니다.
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 토큰이 유효하면 클레임을 반환합니다.
	if claims, ok := refreshToken.Claims.(*model.Claims); ok && refreshToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GetUserByEmail(email string) (*model.User, error) {

	return nil, nil
}

func VerifyPassword(email, password string) bool {
	return true
}
