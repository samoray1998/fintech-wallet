package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct{
	ID primitive.ObjectID  `bson:"_id"`
	FromAccount primitive.ObjectID `bson:"from_account`
	ToAccount primitive.ObjectID `bson:"to_account"`
	Amount float64 `bson:"amount"`
	Fee float64 `bson:"fee"`
	Status string `bson:"status"` // "pending", "completed", "failed"
	CreatedAt time.Time `bson:"created_at"`
}