package service

import (
	"context"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mi-gongan/commention_backend/pkg/db"
	"github.com/mi-gongan/commention_backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCommentsByEmail(email string) ([]model.Comment, error) {
	var comments []model.Comment
	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	cursor, err := collection.Find(context.Background(), bson.M{"owner_email": email})
	if err != nil {
		return comments, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var comment model.Comment
		cursor.Decode(&comment)
		comments = append(comments, comment)
	}
	return comments, nil
}

func CreateComment(comment model.Comment) error {
	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	_, err := collection.InsertOne(context.Background(), comment)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentByID(c *gin.Context, id string) error {
	isOwner, err := checkCommentOwner(c, id)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("you are not the owner of this comment")
	}

	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func UpdateCommentByID(c *gin.Context, id string, content string) error {
	isOwner, err := checkCommentOwner(c, id)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("you are not the owner of this comment")
	}

	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"content": content}})
	if err != nil {
		return err
	}
	return nil
}

func GetCommentByID(c *gin.Context, id string) (*model.Comment, error) {
	isOwner, err := checkCommentOwner(c, id)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("you are not the owner of this comment")
	}

	var comment model.Comment
	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return &comment, err
	}
	return &comment, nil
}

func UpdateCommentIsDisplayedByID(c *gin.Context, id string, isDisplayed bool) error {
	isOwner, err := checkCommentOwner(c, id)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("you are not the owner of this comment")
	}

	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"is_displayed": isDisplayed}})
	if err != nil {
		return err
	}
	return nil
}

func UpdateCommentsOrder(c *gin.Context, comments []model.Comment) error {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return errors.New("authorization header is required")
	}
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return err
	}

	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	for _, comment := range comments {
		if comment.OwnerEmail != claims.UserForJWT.Email {
			return errors.New("you are not the owner of this comment")
		}

		_, err := collection.UpdateOne(context.Background(), bson.M{"_id": comment.ID}, bson.M{"$set": bson.M{"order": comment.Order}})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetMaxOrder(email string) int {
	var comment model.Comment
	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	opts := options.FindOne().SetSort(bson.M{"order": -1})
	err := collection.FindOne(context.Background(), bson.M{"owner_email": email}, opts).Decode(&comment)
	if err != nil {
		return 0
	}
	return comment.Order
}

func checkCommentOwner(c *gin.Context, id string) (bool, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return false, errors.New("authorization header is required")
	}
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return false, err
	}
	var comment model.Comment
	collection := db.DBClient.Database(os.Getenv("DB_NAME")).Collection("comments")
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return false, err
	}
	if comment.OwnerEmail != claims.UserForJWT.Email {
		return false, nil
	}
	return true, nil
}
