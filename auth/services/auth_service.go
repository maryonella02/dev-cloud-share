package services

import (
	"auth/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *mongo.Collection
}

func NewAuthService(client *mongo.Client) *Service {
	db := client.Database("auth").Collection("users")
	return &Service{db}
}

func (s *Service) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	_, err = s.db.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// Clear the password before returning the user object
	user.Password = ""
	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	var user models.User
	err := s.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = comparePasswords(user.Password, password)
	if err != nil {
		return nil, err
	}

	// Clear the password before returning the user object
	user.Password = ""
	return &user, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
