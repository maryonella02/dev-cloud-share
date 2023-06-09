package services

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"resource-manager/config"
	"resource-manager/models"
	"time"
)

type ResourceService struct {
	db *mongo.Database
}

func NewResourceService(db *mongo.Database) *ResourceService {
	return &ResourceService{db}
}

func (rs *ResourceService) RegisterResource(resource *models.Resource) error {
	_, err := rs.db.Collection("resources").InsertOne(context.Background(), resource)
	return err
}

func (rs *ResourceService) GetResources() ([]models.Resource, error) {
	cursor, err := rs.db.Collection("resources").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(cursor, context.Background())

	var resources []models.Resource
	err = cursor.All(context.Background(), &resources)
	return resources, err
}

func (rs *ResourceService) UpdateResource(resourceID string, updatedResource *models.Resource) error {
	id, _ := primitive.ObjectIDFromHex(resourceID)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"cpu_cores": updatedResource.CPUCores,
			"memory_mb": updatedResource.MemoryMB,
			"lender_id": updatedResource.LenderID,
		},
	}

	_, err := rs.db.Collection("resources").UpdateOne(context.Background(), filter, update)
	return err
}

func (rs *ResourceService) DeleteResource(resourceID string) error {
	id, _ := primitive.ObjectIDFromHex(resourceID)
	filter := bson.M{"_id": id}

	_, err := rs.db.Collection("resources").DeleteOne(context.Background(), filter)
	return err
}

var ErrResourceNotFound = errors.New("no suitable resource found")

func (rs *ResourceService) AllocateResource(borrowerID string, request models.Resource) (*models.Resource, error) {
	bID, _ := primitive.ObjectIDFromHex(borrowerID)

	// Find a resource that matches the request
	filter := bson.M{
		"cpu_cores": bson.M{"$gte": request.CPUCores},
		"memory_mb": bson.M{"$gte": request.MemoryMB},
		"$or": []bson.M{
			{"borrower_id": bson.M{"$exists": false}},
			{"borrower_id": primitive.ObjectID{}},
		},
	}
	var resource models.Resource
	err := rs.db.Collection("resources").FindOne(context.Background(), filter).Decode(&resource)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	// Assign the resource to the borrower
	resourceFilter := bson.M{"_id": resource.ID}
	resourceUpdate := bson.M{"$set": bson.M{"borrower_id": bID}}
	_, err = rs.db.Collection("resources").UpdateOne(context.Background(), resourceFilter, resourceUpdate)
	if err != nil {
		return nil, err
	}

	// Create a new ResourceUsage entry
	resourceUsage := &models.ResourceUsage{
		ResourceID: resource.ID,
		BorrowerID: bID,
		StartTime:  time.Now(),
	}

	// Insert the ResourceUsage entry into the database
	_, err = rs.db.Collection("resource_usages").InsertOne(context.Background(), resourceUsage)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (rs *ResourceService) ReleaseResource(resourceID string) (*models.Resource, error) {
	rID, _ := primitive.ObjectIDFromHex(resourceID)

	// Retrieve the resource
	resourceFilter := bson.M{"_id": rID, "borrower_id": bson.M{"$ne": primitive.NilObjectID}}
	var resource models.Resource
	err := rs.db.Collection("resources").FindOne(context.Background(), resourceFilter).Decode(&resource)
	if err != nil {
		return nil, err
	}

	// Retrieve the ResourceUsage entry
	usageFilter := bson.M{"resource_id": rID, "borrower_id": resource.BorrowerID, "end_time": bson.M{"$exists": false}}
	var resourceUsage models.ResourceUsage
	err = rs.db.Collection("resource_usages").FindOne(context.Background(), usageFilter).Decode(&resourceUsage)
	if err != nil {
		return nil, err
	}

	// Calculate the usage duration and cost
	duration := time.Since(resourceUsage.StartTime)
	cost, err := calculateUsageCost(resource, duration)
	if err != nil {
		return nil, err
	}

	// Calculate compensation for the lender
	compensation, err := rs.CalculateCompensation(&resource, duration)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Cost: %f; Compensation: %f", cost, compensation))

	// Update the ResourceUsage entry with the calculated cost, compensation, and end time
	usageUpdate := bson.M{"$set": bson.M{"end_time": time.Now(), "cost": cost, "compensation": compensation}}
	_, err = rs.db.Collection("resource_usages").UpdateOne(context.Background(), usageFilter, usageUpdate)
	if err != nil {
		return nil, err
	}

	// Retrieve the lender associated with the resource
	lenderFilter := bson.M{"_id": resource.LenderID}
	var lender models.Lender
	err = rs.db.Collection("lenders").FindOne(context.Background(), lenderFilter).Decode(&lender)
	if err != nil {
		return nil, err
	}

	// Update the lender's reputation by adding 1
	lender.Reputation++

	// Save the updated lender to the database
	lenderUpdate := bson.M{"$set": bson.M{"reputation": lender.Reputation}}
	_, err = rs.db.Collection("lenders").UpdateOne(context.Background(), lenderFilter, lenderUpdate)
	if err != nil {
		return nil, err
	}

	// Set the borrower_id back to the default empty value
	resourceUpdate := bson.M{"$set": bson.M{"borrower_id": primitive.NilObjectID}}
	_, err = rs.db.Collection("resources").UpdateOne(context.Background(), resourceFilter, resourceUpdate)
	if err != nil {
		return nil, err
	}

	resource.BorrowerID = primitive.NilObjectID
	return &resource, nil
}
func (rs *ResourceService) GetFreeResources() ([]models.Resource, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"borrower_id": bson.M{"$exists": false}},
			{"borrower_id": primitive.ObjectID{}},
		},
	}

	var resources []models.Resource
	cursor, err := rs.db.Collection("resources").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var resource models.Resource
		if err = cursor.Decode(&resource); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

func calculateUsageCost(resource models.Resource, duration time.Duration) (float64, error) {
	// Get pricing config for the resource type
	var pricingConfig config.PricingConfig
	for _, pc := range config.PricingModel {
		pricingConfig = pc
	}

	hours := duration.Hours()
	cost := pricingConfig.PricePerCore*float64(resource.CPUCores) +
		pricingConfig.PricePerMB*float64(resource.MemoryMB)

	return cost * hours, nil
}

func (rs *ResourceService) CalculateCompensation(resource *models.Resource, duration time.Duration) (float64, error) {
	ls := NewLenderService(rs.db)
	reputation, err := ls.GetReputation(resource.LenderID.Hex())
	if err != nil {
		return 0, err
	}

	cost, err := calculateUsageCost(*resource, duration)
	if err != nil {
		return 0, err
	}

	// Calculate reputation bonus
	reputationBonus := 1.0
	if reputation >= 100 {
		reputationBonus = 1.1
	} else if reputation >= 50 {
		reputationBonus = 1.05
	}

	compensation := cost * reputationBonus
	return compensation, nil
}

func (rs *ResourceService) ApplyDiscount(cost float64, discountPercentage float64) float64 {
	return cost * (1 - discountPercentage/100)
}

func (rs *ResourceService) CreateBorrower(borrower models.Borrower) (*models.Borrower, error) {
	borrower.ID = primitive.NewObjectID()
	res, err := rs.db.Collection("borrowers").InsertOne(context.Background(), borrower)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		borrower.ID = oid
	}
	return &borrower, nil
}
func (rs *ResourceService) CreateLender(lender models.Lender) (*models.Lender, error) {
	lender.ID = primitive.NewObjectID()
	res, err := rs.db.Collection("lenders").InsertOne(context.Background(), lender)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		lender.ID = oid
	}
	return &lender, nil
}
