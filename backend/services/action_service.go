// services/action_service.go
package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"riskmgt/models"
	"riskmgt/utils"
)

func ValidateActionCreation(action *models.Action) error {
	if action.Title == "" || action.Description == "" {
		return errors.New("title and description are required")
	}
	return nil
}

func SubmitActionForApproval(ctx context.Context, action *models.Action) error {
	if action.Status != "Draft" {
		return errors.New("only draft actions can be submitted")
	}

	update := bson.M{
		"$set": bson.M{
			"status":    "PendingApproval",
			"updatedAt": time.Now(),
		},
	}
	_, err := ActionCollection.UpdateOne(ctx, bson.M{"_id": action.ID}, update)
	if err != nil {
		return err
	}

	approval := models.Approval{
		ID:              primitive.NewObjectID(),
		TargetID:        action.ID,
		TargetType:      "Action",
		Status:          "Pending",
		RequestedBy:     action.OwnerID,
		RequestedByName: action.OwnerName,
		RequestedAt:     time.Now(),
	}

	_, err = ApprovalCollection.InsertOne(ctx, approval)
	if err != nil {
		return err
	}

	LogAudit(ctx, action.OwnerID, action.OwnerName, "Submitted", "Action", action.ID, nil)

	return nil
}

func PaginateActions(ctx context.Context, baseQuery bson.M, filters bson.M, page utils.Pagination) ([]models.Action, int64, error) {
	for k, v := range filters {
		baseQuery[k] = v
	}

	count, err := ActionCollection.CountDocuments(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((page.Page-1) * page.PageSize)).
		SetLimit(int64(page.PageSize))

	cursor, err := ActionCollection.Find(ctx, baseQuery, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var actions []models.Action
	if err = cursor.All(ctx, &actions); err != nil {
		return nil, 0, err
	}

	return actions, count, nil
}