package model

import "github.com/dgrijalva/jwt-go"

// User 구조체는 사용자 정보를 나타냅니다.
type UserForJWT struct {
	Email string `json:"email"`
}

// Claims 구조체는 JWT 토큰의 클레임을 나타냅니다.
type Claims struct {
	UserForJWT `json:"user"`
	jwt.StandardClaims
}
