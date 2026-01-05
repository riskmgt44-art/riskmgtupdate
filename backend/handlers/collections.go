// handlers/collections.go
package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"

	"riskmgt/services"
)

var (
	riskCollection     *mongo.Collection
	actionCollection   *mongo.Collection
	approvalCollection *mongo.Collection
	auditCollection    *mongo.Collection
	userCollection     *mongo.Collection
)

func InitCollections() {
	riskCollection     = services.RiskCollection
	actionCollection   = services.ActionCollection
	approvalCollection = services.ApprovalCollection
	auditCollection    = services.AuditCollection
	userCollection     = services.UserCollection
}