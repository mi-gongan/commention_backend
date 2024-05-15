package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mi-gongan/commention_backend/pkg/db"
	"github.com/mi-gongan/commention_backend/pkg/route"
)

func init() {
	godotenv.Load(".env")
	db.ConnectDB()
}

func main() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	route.RegisterAuthRoutes(r, "/auth")

	r.Run(":8080")
}
