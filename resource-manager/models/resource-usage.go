package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ResourceUsage struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ResourceID   primitive.ObjectID `bson:"resource_id" json:"resource_id"`
	BorrowerID   primitive.ObjectID `bson:"borrower_id" json:"borrower_id"`
	StartTime    time.Time          `bson:"start_time" json:"start_time"`
	EndTime      *time.Time         `bson:"end_time,omitempty" json:"end_time,omitempty"`
	Cost         float64            `bson:"cost,omitempty" json:"cost,omitempty"`
	Compensation float64            `bson:"compensation,omitempty" json:"compensation,omitempty"`
}
