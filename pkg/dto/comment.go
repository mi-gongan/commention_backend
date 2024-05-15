package dto

import "github.com/mi-gongan/commention_backend/pkg/model"

type GetCommentByIDRequest struct {
	ID string `json:"id" binding:"required"`
}

type CreateCommentRequest struct {
	ToEmail string `json:"toEmail" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateCommentByIDRequest struct {
	ID      string `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateCommentIsDisplayedByIDRequest struct {
	ID          string `json:"id" binding:"required"`
	IsDisplayed bool   `json:"isDisplayed" binding:"required"`
}

type UpdateCommentsOrderRequest struct {
	Comments []model.Comment `json:"comments" binding:"required"`
}
