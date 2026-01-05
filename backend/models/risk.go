// models/risk.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Risk struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title               string             `bson:"title" json:"title"`
	Description         string             `bson:"description" json:"description"`
	Category            string             `bson:"category" json:"category"`
	Likelihood          int                `bson:"likelihood" json:"likelihood"`                           // 1-5
	Impact              int                `bson:"impact" json:"impact"`                                   // 1-5
	ResidualLikelihood  int                `bson:"residualLikelihood" json:"residualLikelihood"`           // after mitigation
	ResidualImpact      int                `bson:"residualImpact" json:"residualImpact"`
	Status              string             `bson:"status" json:"status"`                                   // Draft, PendingApproval, Approved, Rejected, Closed
	OwnerID             primitive.ObjectID `bson:"ownerID" json:"ownerID"`
	OwnerName           string             `bson:"ownerName" json:"ownerName"`
	LinkedActions       []primitive.ObjectID `bson:"linkedActions,omitempty" json:"linkedActions"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
	ApprovedBy          primitive.ObjectID `bson:"approvedBy,omitempty" json:"approvedBy,omitempty"`
	ApprovedAt          *time.Time         `bson:"approvedAt,omitempty" json:"approvedAt,omitempty"`
	RejectedBy          primitive.ObjectID `bson:"rejectedBy,omitempty" json:"rejectedBy,omitempty"`
	RejectedAt          *time.Time         `bson:"rejectedAt,omitempty" json:"rejectedAt,omitempty"`
	RejectionReason     string             `bson:"rejectionReason,omitempty" json:"rejectionReason,omitempty"`
}

type RiskUpdate struct {
	Title              string `json:"title,omitempty"`
	Description        string `json:"description,omitempty"`
	Category           string `json:"category,omitempty"`
	Likelihood         int    `json:"likelihood,omitempty"`
	Impact             int    `json:"impact,omitempty"`
	ResidualLikelihood int    `json:"residualLikelihood,omitempty"`
	ResidualImpact     int    `json:"residualImpact,omitempty"`
}

type RAMUpdate struct {
	ResidualLikelihood int    `json:"residualLikelihood" bson:"residualLikelihood"`
	ResidualImpact     int    `json:"residualImpact" bson:"residualImpact"`
	UpdatedBy          primitive.ObjectID `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updatedAt"`
}