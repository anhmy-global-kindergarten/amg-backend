package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PostId     string             `bson:"post_id" json:"postId"`
	AuthorId   string             `bson:"author_id,omitempty" json:"authorId,omitempty"`
	AuthorName string             `bson:"author_name" json:"authorName"`
	Content    string             `bson:"content" json:"content"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updatedAt"`
}

type CreateCommentPayload struct {
	PostId     string `json:"postId"`
	AuthorId   string `json:"authorId,omitempty"`
	AuthorName string `json:"authorName"`
	Content    string `json:"content"`
}
