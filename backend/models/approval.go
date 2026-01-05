// models/approval.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Approval struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TargetID        primitive.ObjectID `bson:"targetID" json:"targetID"`                  // Risk or Action ID
	TargetType      string             `bson:"targetType" json:"targetType"`              // "Risk" or "Action"
	Status          string             `bson:"status" json:"status"`                      // Pending, Approved, Rejected
	RequestedBy     primitive.ObjectID `bson:"requestedBy" json:"requestedBy"`
	RequestedByName string             `bson:"requestedByName" json:"requestedByName"`
	RequestedAt     time.Time          `bson:"requestedAt" json:"requestedAt"`
	ReviewedBy      primitive.ObjectID `bson:"reviewedBy,omitempty" json:"reviewedBy,omitempty"`
	ReviewedByName  string             `bson:"reviewedByName,omitempty" json:"reviewedByName,omitempty"`
	ReviewedAt      *time.Time         `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"`
	Reason          string             `bson:"reason,omitempty" json:"reason,omitempty"` // for rejection
	Level           string             `bson:"level,omitempty" json:"level,omitempty"`   // e.g., "RiskManager", "Executive" for multi-level
}