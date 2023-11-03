package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dinethpiyumantha/eventeo-gql-service/graph/model"
	"github.com/dinethpiyumantha/eventeo-gql-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	utils.LoadEnv()
	conString := string(os.Getenv("DATABASE_URL"))
	client, err := mongo.NewClient(options.Client().ApplyURI(conString))
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

	user := db.GetUser(eventInfo.OrganizerID)
	inserg, err := eventCollec.InsertOne(ctx, bson.M{"title": eventInfo.Title, "description": eventInfo.Description, "url": eventInfo.URL, "organizer": user})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnEventListing := model.EventListing{ID: insertedID, Title: eventInfo.Title, Organizer: user, Description: eventInfo.Description, URL: eventInfo.URL}
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

func (db *DB) GetEventsPaginated(page int, limit int) []*model.EventListing {
	eventCollec := db.client.Database("eventeo-db").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Calculate the number of items to skip based on the page and limit
	skip := (page - 1) * limit

	// Set options for pagination
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	// Execute the paginated query
	cursor, err := eventCollec.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var eventListings []*model.EventListing
	if err = cursor.All(ctx, &eventListings); err != nil {
		log.Fatal(err)
	}

	return eventListings
}

// User
func (db *DB) CreateUser(userInfo model.CreateUserInput) *model.User {
	userCollec := db.client.Database("eventeo-db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := userCollec.InsertOne(ctx, bson.M{"name": userInfo.Name, "email": userInfo.Email, "password": userInfo.Password, "role": userInfo.Role})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnUser := model.User{ID: insertedID, Name: userInfo.Name, Email: userInfo.Email, Password: userInfo.Password, Role: userInfo.Role}
	return &returnUser
}

func (db *DB) GetUsers() []*model.User {
	userCollec := db.client.Database("eventeo-db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var users []*model.User
	cursor, err := userCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		panic(err)
	}
	return users
}

func (db *DB) GetUser(id string) *model.User {
	userCollec := db.client.Database("eventeo-db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var user model.User
	err := userCollec.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	return &user
}
