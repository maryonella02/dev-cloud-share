package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"resource-manager/models"
)

type LenderService struct {
	db *mongo.Database
}

func NewLenderService(db *mongo.Database) *LenderService {
	return &LenderService{db: db}
}

func (ls *LenderService) GetReputation(lenderID string) (int, error) {
	id, _ := primitive.ObjectIDFromHex(lenderID)
	filter := bson.M{"_id": id}
	var lender models.Lender
	err := ls.db.Collection("lenders").FindOne(context.Background(), filter).Decode(&lender)
	if err != nil {
		return 0, err
	}

	return lender.Reputation, nil
}
