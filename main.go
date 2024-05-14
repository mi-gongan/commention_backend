package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mi-gongan/commention_backend/pkg/route"
)

type User struct {
	Email string `json:"email"`
}

type Claims struct {
	User `json:"user"`
	jwt.StandardClaims
}

func init() {
	godotenv.Load(".env")
}

func main() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	route.RegisterAuthRoutes(r, "/auth")

	r.Run(":8080")
}
