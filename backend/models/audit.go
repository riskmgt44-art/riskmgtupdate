// models/audit.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Audit struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	UserName    string             `bson:"userName" json:"userName"`
	Action      string             `bson:"action" json:"action"` // Created, Updated, Submitted, Approved, Rejected, etc.
	TargetType  string             `bson:"targetType" json:"targetType"` // Risk, Action, User, etc.
	TargetID    primitive.ObjectID `bson:"targetID" json:"targetID"`
	Changes     bson.M             `bson:"changes,omitempty" json:"changes,omitempty"` // diff of before/after
	IPAddress   string             `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	UserAgent   string             `bson:"userAgent,omitempty" json:"userAgent,omitempty"`
}