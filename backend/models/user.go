// models/user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Email         string             `bson:"email" json:"email"`
	PasswordHash  string             `bson:"passwordHash" json:"-"`
	Role          string             `bson:"role" json:"role"` // Analyst, RiskManager, Executive, Viewer, Admin
	Department    string             `bson:"department,omitempty" json:"department,omitempty"`
	IsActive      bool               `bson:"isActive" json:"isActive"`
	LastLogin     *time.Time         `bson:"lastLogin,omitempty" json:"lastLogin,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}