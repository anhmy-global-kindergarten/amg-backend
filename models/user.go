package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"-"` // Không trả về trong JSON
	Name     string             `bson:"name" json:"name"`
	Role     string             `bson:"role" json:"role"` // user / admin / teacher
	CreateAt time.Time          `bson:"create_at" json:"date_created"`
	UpdateAt time.Time          `bson:"update_at" json:"update_at"`
	IsActive bool               `bson:"is_active" json:"is_active"`
}
