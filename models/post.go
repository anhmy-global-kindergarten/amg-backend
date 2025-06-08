package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	Category string             `json:"category" bson:"category"`
	AuthorID primitive.ObjectID `json:"author_id" bson:"author_id"`
	CreateAt time.Time          `json:"create_at" bson:"create_at"`
	UpdateAt time.Time          `json:"update_at" bson:"update_at"`
	Status   string             `json:"status" bson:"status"`
}
