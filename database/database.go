package database

import (
	"context"
	"log"
	"time"

	"github.com/dinethpiyumantha/eventeo-gql-service/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb+srv://dinethpiyumantha:adminpasstTest@testingcluster.g2cxfdf.mongodb.net/eventeo-db?retryWrites=true&w=majority"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetEvent(id string) *model.EventListing {
	eventCollec := db.client.Database("eventeo-db").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var eventListing model.EventListing
	err := eventCollec.FindOne(ctx, filter).Decode(&eventListing)
	if err != nil {
		log.Fatal(err)
	}
	return &eventListing
}

func (db *DB) GetEvents() []*model.EventListing {
	eventCollec := db.client.Database("eventeo-db").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var eventListings []*model.EventListing
	cursor, err := eventCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &eventListings); err != nil {
		panic(err)
	}

	return eventListings
}

func (db *DB) CreateEventListing(eventInfo model.CreateEventListingInput) *model.EventListing {
	eventCollec := db.client.Database("eventeo-db").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := eventCollec.InsertOne(ctx, bson.M{"title": eventInfo.Title, "description": eventInfo.Description, "url": eventInfo.URL, "organizer": eventInfo.Organizer})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnEventListing := model.EventListing{ID: insertedID, Title: eventInfo.Title, Organizer: eventInfo.Organizer, Description: eventInfo.Description, URL: eventInfo.URL}
	return &returnEventListing
}

func (db *DB) UpdateEventListing(eventId string, eventInfo model.UpdateEventListingInput) *model.EventListing {
	eventCollec := db.client.Database("eventeo-db").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateEventInfo := bson.M{}

	if eventInfo.Title != nil {
		updateEventInfo["title"] = eventInfo.Title
	}
	if eventInfo.Description != nil {
		updateEventInfo["description"] = eventInfo.Description
	}
	if eventInfo.URL != nil {
		updateEventInfo["url"] = eventInfo.URL
	}

	_id, _ := primitive.ObjectIDFromHex(eventId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateEventInfo}

	results := eventCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var eventListing model.EventListing

	if err := results.Decode(&eventListing); err != nil {
		log.Fatal(err)
	}

	return &eventListing
}

func (db *DB) DeleteEventListing(eventId string) *model.DeleteEventResponse {
	eventCollec := db.client.Database("eventeo-db").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(eventId)
	filter := bson.M{"_id": _id}
	_, err := eventCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteEventResponse{DeleteEventID: eventId}
}
