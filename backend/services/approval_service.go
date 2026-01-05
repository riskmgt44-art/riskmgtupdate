// services/approval_service.go
package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"riskmgt/models"
)

func ProcessApproval(ctx context.Context, approvalID primitive.ObjectID, decision string, userID primitive.ObjectID, userRole string, reason ...string) error {
	var approval models.Approval
	err := ApprovalCollection.FindOne(ctx, bson.M{"_id": approvalID, "status": "Pending"}).Decode(&approval)
	if err != nil {
		return errors.New("approval not found or already processed")
	}

	updateApproval := bson.M{
		"$set": bson.M{
			"status":         decision,
			"reviewedBy":     userID,
			"reviewedByName": ctx.Value("userName").(string),
			"reviewedAt":     time.Now(),
		},
	}
	if len(reason) > 0 {
		updateApproval["$set"].(bson.M)["reason"] = reason[0]
	}

	_, err = ApprovalCollection.UpdateOne(ctx, bson.M{"_id": approvalID}, updateApproval)
	if err != nil {
		return err
	}

	targetCollection := RiskCollection
	if approval.TargetType == "Action" {
		targetCollection = ActionCollection
	}

	status := "Approved"
	if decision == "Rejected" {
		status = "Rejected"
	}

	updateTarget := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}
	if decision == "Approved" {
		updateTarget["$set"].(bson.M)["approvedBy"] = userID
		updateTarget["$set"].(bson.M)["approvedAt"] = time.Now()
	} else {
		updateTarget["$set"].(bson.M)["rejectedBy"] = userID
		updateTarget["$set"].(bson.M)["rejectedAt"] = time.Now()
		if len(reason) > 0 {
			updateTarget["$set"].(bson.M)["rejectionReason"] = reason[0]
		}
	}

	_, err = targetCollection.UpdateOne(ctx, bson.M{"_id": approval.TargetID}, updateTarget)
	if err != nil {
		return err
	}

	LogAudit(ctx, userID, ctx.Value("userName").(string), decision, approval.TargetType, approval.TargetID, nil)

	return nil
}