package services

import (
	"context"
	"dev-cloud-share/resource-manager/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
