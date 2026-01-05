// services/risk_service.go
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

func ValidateRiskCreation(risk *models.Risk) error {
	if risk.Title == "" || risk.Description == "" || risk.Category == "" {
		return errors.New("title, description and category are required")
	}
	if risk.Likelihood < 1 || risk.Likelihood > 5 || risk.Impact < 1 || risk.Impact > 5 {
		return errors.New("likelihood and impact must be between 1 and 5")
	}
	return nil
}

func SubmitRiskForApproval(ctx context.Context, risk *models.Risk) error {
	if risk.Status != "Draft" {
		return errors.New("only draft risks can be submitted")
	}

	update := bson.M{
		"$set": bson.M{
			"status":    "PendingApproval",
			"updatedAt": time.Now(),
		},
	}
	_, err := RiskCollection.UpdateOne(ctx, bson.M{"_id": risk.ID}, update)
	if err != nil {
		return err
	}

	approval := models.Approval{
		ID:              primitive.NewObjectID(),
		TargetID:        risk.ID,
		TargetType:      "Risk",
		Status:          "Pending",
		RequestedBy:     risk.OwnerID,
		RequestedByName: risk.OwnerName,
		RequestedAt:     time.Now(),
	}

	_, err = ApprovalCollection.InsertOne(ctx, approval)
	if err != nil {
		return err
	}

	LogAudit(ctx, risk.OwnerID, risk.OwnerName, "Submitted", "Risk", risk.ID, nil)

	return nil
}

func UpdateRiskRAM(ctx context.Context, riskID primitive.ObjectID, ramUpdate *models.RAMUpdate, userID primitive.ObjectID) error {
	var risk models.Risk
	err := RiskCollection.FindOne(ctx, bson.M{"_id": riskID}).Decode(&risk)
	if err != nil {
		return errors.New("risk not found")
	}

	if risk.Status != "Approved" {
		return errors.New("RAM can only be updated on approved risks")
	}

	update := bson.M{
		"$set": bson.M{
			"residualLikelihood": ramUpdate.ResidualLikelihood,
			"residualImpact":     ramUpdate.ResidualImpact,
			"updatedAt":          time.Now(),
		},
	}

	_, err = RiskCollection.UpdateOne(ctx, bson.M{"_id": riskID}, update)
	if err != nil {
		return err
	}

	LogAudit(ctx, userID, "Unknown", "UpdatedRAM", "Risk", riskID, ramUpdate)

	return nil
}

func PaginateRisks(ctx context.Context, baseQuery bson.M, filters bson.M, page utils.Pagination) ([]models.Risk, int64, error) {
	for k, v := range filters {
		baseQuery[k] = v
	}

	count, err := RiskCollection.CountDocuments(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((page.Page-1) * page.PageSize)).
		SetLimit(int64(page.PageSize))

	cursor, err := RiskCollection.Find(ctx, baseQuery, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var risks []models.Risk
	if err = cursor.All(ctx, &risks); err != nil {
		return nil, 0, err
	}

	return risks, count, nil
}