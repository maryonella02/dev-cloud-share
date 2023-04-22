package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Resource struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type     string             `bson:"type" json:"type"`
	Capacity int                `bson:"capacity" json:"capacity"`
	LenderID primitive.ObjectID `bson:"lender_id,omitempty" json:"lender_id,omitempty"`
}

type Lender struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	ContactInfo string             `bson:"contact_info" json:"contact_info"`
}

type Borrower struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string               `bson:"name" json:"name"`
	ContactInfo string               `bson:"contact_info" json:"contact_info"`
	Resources   []primitive.ObjectID `bson:"resources,omitempty" json:"resources,omitempty"`
}
