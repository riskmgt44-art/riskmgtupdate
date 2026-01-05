// services/collections.go
package services

import "go.mongodb.org/mongo-driver/mongo"

var (
	RiskCollection     *mongo.Collection
	ActionCollection   *mongo.Collection
	ApprovalCollection *mongo.Collection
	AuditCollection    *mongo.Collection
	UserCollection     *mongo.Collection
)