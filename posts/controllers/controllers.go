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
	"gopkg.in/mgo.v2/bson"
)

// Create a new post. UserId must be sent
func Create(w http.ResponseWriter, r *http.Request) {
	// Create a object where we are going to store the post values
	post, err := decodePostFromBody(w, r)
	if err != nil {
		return
	}

	// Fetch posts collection
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
	common.SendResponse(w, http.StatusOK, []byte(data), "")
}

// Fetch one post by id
func GetById(w http.ResponseWriter, r *http.Request) {
	// Fetch post_id from request URL
	oid, err := decodeIdFromURL(w, r)
	if err != nil {
		return
	}

	// Create a object where we are going to store the fetched data
	var post models.Post

	// Fetch posts collection
	var postsCollection *mongo.Collection
	postsCollection, ctx := fetchConnection()

	// Fetch the post that matches with the provided ObjectId
	filter := bson.M{"_id": oid}
	err = postsCollection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		data := `{
			"message": "Couldnt fetch any post with provided 'id'"
		}`
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return
	}

	// Return the fetched post as a json
	json, err := json.Marshal(post)

	data := fmt.Sprintf(`{
		"message": "Post fetched successfully",
		"post": %s
	}`, json)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data), "")
}

// Fetch all posts from one user
func GetAllFromUserId(w http.ResponseWriter, r *http.Request) {
	// Fetch user_id from request URL
	oid, err := decodeIdFromURL(w, r)
	if err != nil {
		return
	}

	// Create an array of objects where we are going to store the fetched data
	var posts []models.Post

	// Fetch posts collection
	var postsCollection *mongo.Collection
	postsCollection, ctx := fetchConnection()

	// Fetch all posts that matches with the provided user ObjectId
	filter := bson.M{"user_id": oid}
	cursor, err := postsCollection.Find(ctx, filter)
	if err != nil {
		data := `{
			"message": "We couldnt fetch any post with provided 'id'"
		}`
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Post
		if err = cursor.Decode(&p); err != nil {
			data := `{
				"message": "There was an error decoding fetched posts"
			}`
			common.SendError(w, http.StatusBadRequest, []byte(data))
			return
		}
		posts = append(posts, p)
	}

	// Return the fetched post as a json
	json, err := json.Marshal(posts)

	data := fmt.Sprintf(`{
		"message": "Posts fetched successfully",
		"posts": %s
	}`, json)

	if len(posts) == 0 {
		data = fmt.Sprintf(`{
			"message": "There are 0 post vinculated with provided 'user_id'"
		}`)
	}

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data), "")
}

// Delete the post that matches with provided ObjectId
func Delete(w http.ResponseWriter, r *http.Request) {
	// Fetch post_id from request URL
	oid, err := decodeIdFromURL(w, r)
	if err != nil {
		return
	}

	// Fetch posts collection
	var postsCollection *mongo.Collection
	postsCollection, ctx := fetchConnection()

	// Delete the post that matches with the provided ObjectId
	filter := bson.M{"_id": oid}
	deleteResult, err := postsCollection.DeleteOne(ctx, filter)

	if err != nil {
		data := `{
			"message": "There was an error deleting the post"
		}`
		common.SendError(w, http.StatusInternalServerError, []byte(data))
		return
	}

	if deleteResult.DeletedCount == 0 {
		data := `{
			"message": "There weren't any post that match with the provided 'id'"
		}`
		common.SendError(w, http.StatusOK, []byte(data))
		return
	}

	data := fmt.Sprintf(`{
		"message": "Post deleted successfully"
	}`)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data), "")
}

// Update title and description from post that matches de provided id
func Update(w http.ResponseWriter, r *http.Request) {
	// Fetch post_id from request URL
	oid, err := decodeIdFromURL(w, r)
	if err != nil {
		return
	}

	// Create a object where we are going to store the fetched data
	var post models.Post
	post.Id = oid

	// Fetch the title and description from request body
	post, err = decodePostFromBody(w, r)
	if err != nil {
		return
	}

	// Fetch posts collection
	var postsCollection *mongo.Collection
	postsCollection, ctx := fetchConnection()

	// Update the post that matches with the provided ObjectId
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"title":       post.Title,
			"description": post.Description,
		},
	}
	updateResult, err := postsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		data := `{
			"message": "There was an error while trying to update the post"
		}`
		common.SendError(w, http.StatusInternalServerError, []byte(data))
		return
	}

	if updateResult.ModifiedCount == 0 {
		data := `{
			"message": "Couldnt fetch any post with the provided ObjectId"
		}`
		common.SendResponse(w, http.StatusOK, []byte(data), "")
		return
	}

	data := fmt.Sprintf(`{
		"message": "Post updated successfully"
	}`)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data), "")
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

	// On update operation we are not going to send user id
	if data.UserId != "" {
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
	}

	post.Title = data.Title
	post.Description = data.Description

	return post, nil
}

// Decode id from request body and returns a primitive.ObjectID
func decodeIdFromURL(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, error) {
	urlParam := r.URL.Query().Get("id")

	if urlParam == "" {
		urlParam = r.URL.Query().Get("user_id")
	}

	// It's a string, should turn it into a primitive.ID
	oid, err := primitive.ObjectIDFromHex(urlParam)
	if err != nil {
		data := fmt.Sprintf(`{
			"message": "The provided 'id' is invalid. Try sending a valid one",
			"error": %s
		}`, err.Error())
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return oid, err
	}

	return oid, nil
}

// Fetch posts collection and the context
func fetchConnection() (*mongo.Collection, context.Context) {
	// Fetch the database
	mongoClient := database.GetMongoClient()
	postsCollection := mongoClient.Database("database").Collection("posts")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return postsCollection, ctx
}
