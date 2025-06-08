package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Candidate struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StudentName string             `json:"student_name" bson:"student_name"`
	Gender      string             `json:"gender" bson:"gender"`
	DateOfBirth string             `json:"dob" bson:"dob"`
	ParentName  string             `json:"parent_name" bson:"parent_name"`
	Address     string             `json:"address" bson:"address"`
	Phone       string             `json:"phone" bson:"phone"`
	Status      string             `json:"status" bson:"status"`
	CreateAt    time.Time          `json:"create_at" bson:"create_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
}
