// services/audit_service.go
package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"riskmgt/models"
)

func LogAudit(ctx context.Context, userID primitive.ObjectID, userName, action, targetType string, targetID primitive.ObjectID, changes interface{}) {
	audit := models.Audit{
		ID:         primitive.NewObjectID(),
		Timestamp:  time.Now(),
		UserID:     userID,
		UserName:   userName,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
	}

	if changes != nil {
		if m, ok := changes.(bson.M); ok {
			audit.Changes = m
		}
	}

	if r := ctx.Value("remoteAddr"); r != nil {
		if addr, ok := r.(string); ok {
			audit.IPAddress = addr
		}
	}
	if r := ctx.Value("userAgent"); r != nil {
		if ua, ok := r.(string); ok {
			audit.UserAgent = ua
		}
	}

	AuditCollection.InsertOne(context.Background(), audit)
}