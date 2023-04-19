package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define the data structures for resources and requests
type Resource struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Type        string             `bson:"type"`
	Description string             `bson:"description"`
	Lender      string             `bson:"lender"`
	TotalUnits  int                `bson:"total_units"`
	UsedUnits   int                `bson:"used_units"`
}

type Request struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ResourceType string             `bson:"resource_type"`
	Borrower     string             `bson:"borrower"`
	Units        int                `bson:"units"`
	Status       string             `bson:"status"`
}
