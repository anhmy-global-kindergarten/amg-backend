package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ImageStatus string

const (
	ImageStatusPending ImageStatus = "pending"
	ImageStatusUsed    ImageStatus = "used"
)

type UploadedImage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Filename  string             `bson:"filename" json:"filename"`
	Path      string             `bson:"path" json:"path"`
	URL       string             `bson:"url" json:"url"`
	Status    ImageStatus        `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
