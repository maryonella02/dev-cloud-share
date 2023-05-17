package services

import (
	"auth/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db *mongo.Collection
}

func NewAuthService(client *mongo.Client) *Service {
	db := client.Database("auth").Collection("users")
	return &Service{db}
}

func (s *Service) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := s.db.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User
	err := s.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
