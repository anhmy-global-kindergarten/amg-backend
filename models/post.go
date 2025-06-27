package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"`
	HeaderImage string             `json:"header_image" bson:"header_image"`
	Category    string             `json:"category" bson:"category"`
	Author      string             `json:"author" bson:"author"`
	CreateAt    time.Time          `json:"create_at" bson:"create_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	Status      string             `json:"status" bson:"status"`
}

type PostDetailResponse struct {
	Post   Post            `json:"post"`
	Images []UploadedImage `json:"images"`
}

type LandingPageContent struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Key       string                 `bson:"key" json:"key"`
	Content   map[string]interface{} `bson:"content" json:"content"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
}
