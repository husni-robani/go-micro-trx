package data

import "go.mongodb.org/mongo-driver/v2/mongo"

type Models struct {
	Log Log
}

type Log struct {

}

var client *mongo.Client

func New(mongoClient *mongo.Client) Models {
	client = mongoClient

	return Models{}
}