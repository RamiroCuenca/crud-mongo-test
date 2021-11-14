package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for Post object
type Post struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Title       string             `json:"username,omitempty" bson:"username,omitempty"`
	Description string             `json:"password,omitempty" bson:"password,omitempty"`
}
