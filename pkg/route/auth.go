package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mi-gongan/commention_backend/pkg/handler"
)

func RegisterAuthRoutes(baseRouter gin.IRouter, relativeUrl string) {
	router := baseRouter.Group(relativeUrl)
	router.GET("/", handler.GetAuthHandler)
	router.POST("/login", handler.LoginHandler)
	router.POST("/refresh", handler.RefreshTokenHandler)
}
