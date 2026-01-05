// models/action.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Action struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	RiskID          primitive.ObjectID `bson:"riskID" json:"riskID"`                              // linked risk
	DueDate         *time.Time         `bson:"dueDate,omitempty" json:"dueDate,omitempty"`
	AssigneeID      primitive.ObjectID `bson:"assigneeID,omitempty" json:"assigneeID,omitempty"`
	AssigneeName    string             `bson:"assigneeName,omitempty" json:"assigneeName,omitempty"`
	Status          string             `bson:"status" json:"status"`                              // Draft, PendingApproval, Approved, InProgress, Completed, Overdue
	OwnerID         primitive.ObjectID `bson:"ownerID" json:"ownerID"`
	OwnerName       string             `bson:"ownerName" json:"ownerName"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	ApprovedBy      primitive.ObjectID `bson:"approvedBy,omitempty" json:"approvedBy,omitempty"`
	ApprovedAt      *time.Time         `bson:"approvedAt,omitempty" json:"approvedAt,omitempty"`
	RejectedBy      primitive.ObjectID `bson:"rejectedBy,omitempty" json:"rejectedBy,omitempty"`
	RejectedAt      *time.Time         `bson:"rejectedAt,omitempty" json:"rejectedAt,omitempty"`
	RejectionReason string             `bson:"rejectionReason,omitempty" json:"rejectionReason,omitempty"`
	CompletedAt     *time.Time         `bson:"completedAt,omitempty" json:"completedAt,omitempty"`
}

type ActionUpdate struct {
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	DueDate      *time.Time         `json:"dueDate,omitempty"`
	AssigneeID   primitive.ObjectID `json:"assigneeID,omitempty"`
	AssigneeName string             `json:"assigneeName,omitempty"`
}