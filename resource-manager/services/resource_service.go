package services

import (
	"context"
	"dev-cloud-share/resource-manager/config"
	"dev-cloud-share/resource-manager/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ResourceService struct {
	db *mongo.Database
}

type ResourceRequest struct {
	ResourceType string `json:"resource_type"`
	MinCPUCores  int    `json:"min_cpu_cores"`
	MinMemoryMB  int    `json:"min_memory_mb"`
	MinStorageGB int    `json:"min_storage_gb"`
}

// TODO: extend this struct further to include other factors such as availability, location, or cost

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
			"type":       updatedResource.Type,
			"cpu_cores":  updatedResource.CPUCores,
			"memory_mb":  updatedResource.MemoryMB,
			"storage_gb": updatedResource.StorageGB,
			"lender_id":  updatedResource.LenderID,
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

func (rs *ResourceService) AllocateResource(borrowerID string, request ResourceRequest) (*models.Resource, error) {
	bID, _ := primitive.ObjectIDFromHex(borrowerID)

	// Find a resource that matches the request
	filter := bson.M{
		"type":       request.ResourceType,
		"cpu_cores":  bson.M{"$gte": request.MinCPUCores},
		"memory_mb":  bson.M{"$gte": request.MinMemoryMB},
		"storage_gb": bson.M{"$gte": request.MinStorageGB},
		"lender_id": bson.M{
			"$exists": false,
		},
	}
	var resource models.Resource
	err := rs.db.Collection("resources").FindOne(context.Background(), filter).Decode(&resource)
	if err != nil {
		return nil, err
	}

	// Assign the resource to the borrower
	borrowerFilter := bson.M{"_id": bID}
	borrowerUpdate := bson.M{"$push": bson.M{"resources": resource.ID}}
	_, err = rs.db.Collection("borrowers").UpdateOne(context.Background(), borrowerFilter, borrowerUpdate)
	if err != nil {
		return nil, err
	}

	// Set the 'LenderID' field of the resource to the borrower's ID
	resourceFilter := bson.M{"_id": resource.ID}
	resourceUpdate := bson.M{"$set": bson.M{"lender_id": bID}}
	_, err = rs.db.Collection("resources").UpdateOne(context.Background(), resourceFilter, resourceUpdate)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (rs *ResourceService) ReleaseResource(allocationID string) error {
	// Assuming allocationID is actually the resourceID
	resourceID, _ := primitive.ObjectIDFromHex(allocationID)

	// Find the resource, get the borrower ID, and remove the resource from the borrower's 'Resources' field
	resourceFilter := bson.M{"_id": resourceID}
	var resource models.Resource
	err := rs.db.Collection("resources").FindOne(context.Background(), resourceFilter).Decode(&resource)
	if err != nil {
		return err
	}

	borrowerFilter := bson.M{"_id": resource.LenderID}
	borrowerUpdate := bson.M{"$pull": bson.M{"resources": resource.ID}}
	_, err = rs.db.Collection("borrowers").UpdateOne(context.Background(), borrowerFilter, borrowerUpdate)
	if err != nil {
		return err
	}

	// Clear the 'LenderID' field of the resource
	resourceUpdate := bson.M{"$set": bson.M{"lender_id": nil}}
	_, err = rs.db.Collection("resources").UpdateOne(context.Background(), resourceFilter, resourceUpdate)
	return err
}

func (rs *ResourceService) CalculateCost(resource *models.Resource, duration time.Duration) (float64, error) {
	// Get pricing config for the resource type
	var pricingConfig config.PricingConfig
	for _, pc := range config.PricingModel {
		if pc.ResourceType == resource.Type {
			pricingConfig = pc
			break
		}
	}

	hours := duration.Hours()
	cost := pricingConfig.PricePerCore*float64(resource.CPUCores) +
		pricingConfig.PricePerMB*float64(resource.MemoryMB) +
		pricingConfig.PricePerGB*float64(resource.StorageGB)

	return cost * hours, nil
}

func (rs *ResourceService) CalculateCompensation(resource *models.Resource, duration time.Duration) (float64, error) {
	ls := NewLenderService(rs.db)
	reputation, err := ls.GetReputation(resource.LenderID.Hex())
	if err != nil {
		return 0, err
	}

	cost, err := rs.CalculateCost(resource, duration)
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
