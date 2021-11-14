package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RamiroCuenca/crud-mongo-test/common"
	"github.com/RamiroCuenca/crud-mongo-test/database"
	"github.com/RamiroCuenca/crud-mongo-test/users/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Register a new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	// Create a object where we are going to store the user values
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		return
	}

	// Fetch users collection
	var usersCollection *mongo.Collection
	usersCollection, ctx := fetchConnection()

	// Insert the user on the Collection
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		data := fmt.Sprintf(`{
				"message": "Couldn't register the user",
				"error": %s
			}`, err.Error())
		common.SendError(w, http.StatusInternalServerError, []byte(data))
		return

	}

	user.Id = result.InsertedID.(primitive.ObjectID)

	// Return the created user as a json
	json, err := json.Marshal(user)

	// Generate the JWT token
	// token, err := auth.GenerateToken(user)
	// if err != nil {
	// 	data := fmt.Sprintf(`{
	// 		"message": "User created successfully",
	// 		"user": %s,
	// 		"jwt": "There was an error generating the JWT token, try loggin in"
	// 	}`, json)
	// 	common.SendError(w, http.StatusOK, []byte(data))
	// 	return
	// }

	// data := fmt.Sprintf(`{
	// 	"message": "User created successfully",
	// 	"user": %s,
	// 	"jwt": %s
	// }`, json, token)

	// Send response
	// common.SendResponse(w, http.StatusOK, []byte(data))

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
	user, err := decodeUserFromBody(w, r)
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

	// Generate the JWT token
	// token, err := auth.GenerateToken(user)
	// if err != nil {
	// 	data := fmt.Sprintf(`{
	// 		"message": "User logged in successfully",
	// 		"user": %s,
	// 		"jwt": "There was an error generating the JWT token, try loggin in again"
	// 	}`, json)
	// 	common.SendError(w, http.StatusOK, []byte(data))
	// 	return
	// }

	// data := fmt.Sprintf(`{
	// 	"message": "User logged in successfully",
	// 	"user": %s,
	// 	"jwt": %s
	// }`, json, token)

	// // Send response
	// common.SendResponse(w, http.StatusOK, []byte(data))

	data := fmt.Sprintf(`{
		"message": "User logged in successfully",
		"user": %s
	}`, json)

	// Send response
	common.SendResponse(w, http.StatusOK, []byte(data))
}

// Delete a user that matches with the provided ObjectId
func Delete(w http.ResponseWriter, r *http.Request) {
	// Fetch the id from request URL
	oid, err := decodeIdFromURL(w, r)
	if err != nil {
		return
	}

	// Fetch users collection
	var usersCollection *mongo.Collection
	usersCollection, ctx := fetchConnection()

	// Delete the user that matches with the provided ObjectId
	filter := bson.M{"_id": oid}
	result, err := usersCollection.DeleteOne(ctx, filter)
	if err != nil {
		data := `{
			"message": "There was an error while trying to delete the user"
		}`
		common.SendError(w, http.StatusInternalServerError, []byte(data))
		return
	}

	if result.DeletedCount == 0 {
		data := `{
			"message": "Couldnt fetch any object with the provided ObjectId"
		}`
		common.SendResponse(w, http.StatusOK, []byte(data))
		return
	}

	data := `{
		"message": "User deleted successfully"
	}`
	common.SendResponse(w, http.StatusOK, []byte(data))
}

func Update(w http.ResponseWriter, r *http.Request) {
	// Fetch the id from request URL
	oid, err := decodeIdFromURL(w, r)
	if err != nil {
		return
	}

	// Fetch the password from request Body
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		return
	}

	// Fetch users collection
	var usersCollection *mongo.Collection
	usersCollection, ctx := fetchConnection()

	// Update the password from the user that matches with the provided ObjectId
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"password": user.Password,
		},
	}

	updateResult, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		data := `{
			"message": "There was an error while trying to update the user"
		}`
		common.SendError(w, http.StatusInternalServerError, []byte(data))
		return
	}

	if updateResult.ModifiedCount == 0 {
		data := `{
			"message": "Couldnt fetch any object with the provided ObjectId"
		}`
		common.SendResponse(w, http.StatusOK, []byte(data))
		return
	}

	data := fmt.Sprintf(`{
		"message": "User password updated successfully"
	}`)
	common.SendResponse(w, http.StatusOK, []byte(data))
}

// Decodes username and password from request body and returns a User object
func decodeUserFromBody(w http.ResponseWriter, r *http.Request) (models.User, error) {
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

// Decode id from request body and returns a primitive.ObjectID
func decodeIdFromURL(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, error) {
	urlParam := r.URL.Query().Get("id")

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

// Fetch users collection and the context
func fetchConnection() (*mongo.Collection, context.Context) {
	// Fetch the database
	mongoClient := database.GetMongoClient()
	usersCollection := mongoClient.Database("database").Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return usersCollection, ctx
}
