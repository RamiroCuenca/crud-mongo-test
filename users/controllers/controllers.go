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
)

// SignUp or register a new user
func Create(w http.ResponseWriter, r *http.Request) {
	// Create a object where we are going to store the user values
	var user models.User

	// Decode the username and password from request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data := fmt.Sprintf(`{
			"message": "Couldn't decode the username or password from request body",
			"error": %s
		}`, err.Error())
		common.SendError(w, http.StatusBadRequest, []byte(data))
		return
	}

	// Fetch the database
	mongoClient := database.GetMongoClient()
	usersCollection := mongoClient.Database("database").Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

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

	// if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
	// 	user.Id = oid
	// }

	// user.Id = result.InsertedID

	// Return the created user as a json
	json, err := json.Marshal(user)

	data := fmt.Sprintf(`{
		"message": "User created successfully",
		"user": %s
	}`, json)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data))
}
