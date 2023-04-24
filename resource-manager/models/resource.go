package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Resource struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type          string             `bson:"type" json:"type"`
	CPUCores      int                `bson:"cpu_cores" json:"cpu_cores"`
	MemoryMB      int                `bson:"memory_mb" json:"memory_mb"`
	StorageGB     int                `bson:"storage_gb" json:"storage_gb"`
	UsageDuration time.Duration
	LenderID      primitive.ObjectID `bson:"lender_id,omitempty" json:"lender_id,omitempty"`
}

type Borrower struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string               `bson:"name" json:"name"`
	ContactInfo string               `bson:"contact_info" json:"contact_info"`
	Resources   []primitive.ObjectID `bson:"resources,omitempty" json:"resources,omitempty"`
}
