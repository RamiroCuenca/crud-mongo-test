package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RamiroCuenca/crud-mongo-test/common"
	"github.com/RamiroCuenca/crud-mongo-test/database"
	"github.com/RamiroCuenca/crud-mongo-test/posts/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new post. UserId must be sent
func Create(w http.ResponseWriter, r *http.Request) {
	// Create a object where we are going to store the post values
	post, err := decodePostFromBody(w, r)
	if err != nil {
		return
	}

	var postsCollection *mongo.Collection
	postsCollection, ctx := fetchConnection()

	// Insert the post on the collection
	result, err := postsCollection.InsertOne(ctx, post)
	if err != nil {
		data := fmt.Sprintf(`{
				"message": "Couldn't create the post",
				"error": %s
			}`, err.Error())
		common.SendError(w, http.StatusInternalServerError, []byte(data))
		return
	}

	post.Id = result.InsertedID.(primitive.ObjectID)

	// Convert post into json
	json, err := json.Marshal(post)

	data := fmt.Sprintf(`{
		"message": "Post created successfully",
		"post": %s
	}`, json)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data))
}

// Decodes UserId, Title and Description from request body and returns a Post object
func decodePostFromBody(w http.ResponseWriter, r *http.Request) (models.Post, error) {
	var post models.Post

	type receivedData struct {
		UserId      string `json:"user_id,omitempty" bson:"user_id,omitempty"`
		Title       string `json:"title,omitempty" bson:"title,omitempty"`
		Description string `json:"description,omitempty" bson:"description,omitempty"`
	}

	var data receivedData

	// Decode the UserId, Title and Description from request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		data := fmt.Sprintf(`{
		"message": "Couldn't decode the 'user_id', 'title' or 'description' from request body",
		"error": %s
	}`, err.Error())
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return post, err
	}

	// Convert the UserID into a primitive.ObjectID
	post.UserId, err = primitive.ObjectIDFromHex(data.UserId)
	if err != nil {
		data := fmt.Sprintf(`{
		"message": "Provided 'user_id' is invalid. Try sending a valid one.",
		"error": %s
	}`, err.Error())
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return post, err
	}

	post.Title = data.Title
	post.Description = data.Description

	return post, nil
}

// Fetch posts collection and the context
func fetchConnection() (*mongo.Collection, context.Context) {
	// Fetch the database
	mongoClient := database.GetMongoClient()
	postsCollection := mongoClient.Database("database").Collection("posts")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return postsCollection, ctx
}
