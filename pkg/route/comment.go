package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mi-gongan/commention_backend/pkg/handler"
)

func RegisterCommentRoutes(baseRouter gin.IRouter, relativeUrl string) {
	router := baseRouter.Group(relativeUrl)
	router.GET("/", handler.GetComments)
	router.GET("/:id", handler.GetCommentByID)
	router.POST("/", handler.CreateComment)
	router.PUT("/:id", handler.UpdateCommentByID)
	router.DELETE("/:id", handler.DeleteCommentByID)
	router.PATCH("/:id/display", handler.UpdateCommentDisplayByID)
	router.PATCH("/order", handler.UpdateCommentsOrder)

}
