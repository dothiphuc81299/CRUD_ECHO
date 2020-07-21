package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc      string             `json:"desc,omitempty" bson:"desc,omitempty"`
	Completed bool               `json:"completed,omitempty" bson:"completed,omitempty"`
}

