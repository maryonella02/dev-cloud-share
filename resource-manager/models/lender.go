package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lender struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	ContactInfo string             `bson:"contact_info" json:"contact_info"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Reputation  int                `bson:"reputation" json:"reputation"`
}
