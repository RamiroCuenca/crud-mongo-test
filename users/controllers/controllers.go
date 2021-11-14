package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RamiroCuenca/vozy-test/common"
	"github.com/RamiroCuenca/vozy-test/database"
	"github.com/RamiroCuenca/vozy-test/users/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Register a new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	// Create a object where we are going to store the user values
	user, err := decodeDataFromBody(w, r)
	if err != nil {
		return
	}

	// Fetch users collection
	var usersCollection *mongo.Collection
	usersCollection, ctx := fetchConnection()

	// Insert the user on the Collection
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		if err != nil {
			data := fmt.Sprintf(`{
				"message": "Couldn't register the user",
				"error": %s
			}`, err.Error())
			common.SendError(w, http.StatusInternalServerError, []byte(data))
			return
		}
	}

	user.Id = result.InsertedID.(primitive.ObjectID)

	// Return the created user as a json
	json, err := json.Marshal(user)

	data := fmt.Sprintf(`{
		"message": "User created successfully",
		"user": %s
	}`, json)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data))
}

// Log in with an existing user
func SignIn(w http.ResponseWriter, r *http.Request) {
	// Create a object where we are going to store the user values
	user, err := decodeDataFromBody(w, r)
	if err != nil {
		return
	}

	// Check Username and password are not empty
	if user.Username == "" || user.Password == "" {
		data := `{
			"message": "Username nor password can be empty. Try sending valid values"
		}`
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return
	}

	// Fetch users collection
	var usersCollection *mongo.Collection
	usersCollection, ctx := fetchConnection()

	// Search for the user
	err = usersCollection.FindOne(ctx, bson.M{"username": user.Username, "password": user.Password}).Decode(&user)
	if err != nil {
		data := `{
			"message": "We couldnt fetch any user with provided credentials"
		}`
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return
	}

	// Return the fetched user as a json
	json, err := json.Marshal(user)

	data := fmt.Sprintf(`{
		"message": "User logged in successfully",
		"user": %s
	}`, json)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data))
}

// Decodes username and password from request body and returns a User object
func decodeDataFromBody(w http.ResponseWriter, r *http.Request) (models.User, error) {
	var u models.User

	// Decode the username and password from request body
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		data := fmt.Sprintf(`{
		"message": "Couldn't decode the username or password from request body",
		"error": %s
	}`, err.Error())
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return u, err
	}

	return u, nil
}

// Fetch users collection and the context
func fetchConnection() (*mongo.Collection, context.Context) {
	// Fetch the database
	mongoClient := database.GetMongoClient()
	usersCollection := mongoClient.Database("database").Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return usersCollection, ctx
}
