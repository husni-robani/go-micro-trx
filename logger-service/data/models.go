package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Models struct {
	Log LogEntry
}

type LogEntry struct {
	ID int `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	Data string `bson:"data,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

var client *mongo.Client

func New(mongoClient *mongo.Client) Models {
	client = mongoClient

	return Models{}
}

func (l LogEntry) InsertOne(newLog LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	newLog.CreatedAt = time.Now()
	newLog.UpdatedAt = time.Now()

	db := client.Database("mongo-trx").Collection("logs")
	_, err := db.InsertOne(ctx, newLog)
	if err != nil {
		log.Println("Failed to insert log: ", err)
		return err
	}
	
	return nil
}