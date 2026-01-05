// database/db.go
package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"riskmgt/config"
	"riskmgt/services"
)

var (
	Client   *mongo.Client
	Database *mongo.Database

	// Collections – will be initialized after successful connection
	RiskCollection     *mongo.Collection
	ActionCollection   *mongo.Collection
	ApprovalCollection *mongo.Collection
	AuditCollection    *mongo.Collection
	UserCollection     *mongo.Collection
)

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.MongoURI)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Verify connection
	if err = Client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	Database = Client.Database(config.DBName)
	log.Println("Connected to MongoDB Atlas:", config.DBName)

	// Initialize all collections now that Database is ready
	RiskCollection     = Database.Collection("risks")
	ActionCollection   = Database.Collection("actions")
	ApprovalCollection = Database.Collection("approvals")
	AuditCollection    = Database.Collection("audits")
	UserCollection     = Database.Collection("users")

	// Export the initialized collections to the services package
	services.RiskCollection     = RiskCollection
	services.ActionCollection   = ActionCollection
	services.ApprovalCollection = ApprovalCollection
	services.AuditCollection    = AuditCollection
	services.UserCollection     = UserCollection

	return nil
}

func Disconnect() {
	if Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Println("Error disconnecting from MongoDB:", err)
	} else {
		log.Println("Disconnected from MongoDB")
	}
}

// GetCollection is kept for possible future use but is no longer needed for the main collections
func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}