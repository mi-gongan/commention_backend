package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mi-gongan/commention_backend/pkg/db"
	"github.com/mi-gongan/commention_backend/pkg/dto"
	"github.com/mi-gongan/commention_backend/pkg/model"
	"github.com/mi-gongan/commention_backend/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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
	c.JSON(http.StatusOK, gin.H{"user": claims.UserForJWT})
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
	token, err := service.CreateToken(claims.UserForJWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 생성된 새로운 토큰 반환
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SignInHandler(c *gin.Context) {

	// 클라이언트로부터 요청 받기
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := db.ConnectDB()

	// 사용자 정보 조회
	var result bson.M = bson.M{}

	err := client.Database("commention").Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&result)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// 비밀번호 검증
	err = bcrypt.CompareHashAndPassword([]byte(result["password"].(string)), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// 사용자 정보 생성
	user := model.UserForJWT{
		Email: req.Email,
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

func SignUpHandler(c *gin.Context) {
	// 클라이언트로부터 요청 받기
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	client := db.ConnectDB()

	// 이미 존재하는 사용자인지 확인
	var result bson.M = bson.M{}

	err = client.Database("commention").Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&result)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	var row = bson.M{
		"email":    req.Email,
		"password": string(hashedPassword),
		"name":     req.Name,
	}

	_, err = client.Database("commention").Collection("users").InsertOne(c, row)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client.Disconnect(c)

	// 사용자 정보 생성
	user := model.UserForJWT{
		Email: req.Email,
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
