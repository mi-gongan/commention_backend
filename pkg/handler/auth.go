package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mi-gongan/commention_backend/pkg/dto"
	"github.com/mi-gongan/commention_backend/pkg/model"
	"github.com/mi-gongan/commention_backend/pkg/service"
)

func GetAuthHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Verify the token
	claims, err := service.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return the user information
	c.JSON(http.StatusOK, gin.H{"user": claims.User})
}

func RefreshTokenHandler(c *gin.Context) {
	// 클라이언트로부터 요청 받기
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 리프레시 토큰 검증
	claims, err := service.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 새로운 JWT 토큰 생성
	token, err := service.CreateToken(claims.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 생성된 새로운 토큰 반환
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(c *gin.Context) {

	// 클라이언트로부터 요청 받기
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//TODO: Supabase로부터 사용자 정보 가져와서 없으면 return

	//TODO: 가져온 데이터에 대해서 비밀번호 검증

	// 사용자 정보 생성
	user := model.User{
		Email: req.Email,
		// 다른 사용자 정보 필드도 필요에 따라 여기에 추가
	}

	// Create JWT token and refresh token
	token, refreshToken, err := service.CreateTokenWithRefresh(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 생성된 토큰 반환
	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refreshToken})
}
