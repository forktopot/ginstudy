package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Telephone string             `bson:"telephone,omitempty"`
	Password  string             `bson:"password,omitempty"`
}
