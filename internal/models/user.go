package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id`
	Name      string             `bson:"name`
	Password  string             `bson:"password_hash`
	Email     string             `bson:"email`
	KYCStatus string             `bson:"kyc_status"` // "unverified", "pending", "verified"
	UpdatedAt time.Time          `bson:"updated_at`
	CreatedAt time.Time          `bson:"created_at`
}
