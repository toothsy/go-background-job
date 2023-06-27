package dbrepo

import (
	"context"
	"encoding/json"
	"github/toothsy/go-background-job/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// inserts the user parameter in to the database

func (dbRepo *mongoDBRepo) SignUp(user *models.UserPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := dbRepo.AppConfig.MongoDatabase.Collection("userData").InsertOne(ctx, user)
	if err != nil {
		log.Println("\n\n could not insert, err is ", err)
	}

}

// FetchUser fetches the user based on the email filter
func (dbRepo *mongoDBRepo) FetchUser(user *models.UserPayload) (*models.UserPayload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"email": user.Email}

	foundUser := &models.UserPayload{}
	err := dbRepo.AppConfig.MongoDatabase.Collection("userData").FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		// Handle the error, such as returning a custom error or logging it
		log.Println("user does not exist")
		return nil, err
	}
	// log.Println("Fetched User ", foundUser)
	return foundUser, nil
}

// marks the user verified if the generated uuid matches one sent via get request
func (dbRepo *mongoDBRepo) UpdateUserVerification(user *models.UserPayload) {
	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": bson.M{"isverified": true}}
	result, err := dbRepo.AppConfig.MongoDatabase.Collection("userData").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Did not update User, as user not found")
	}
	log.Println("Updated User :", result.MatchedCount)

}

func (dbRepo *mongoDBRepo) SearchUserImage(email string) ([]json.RawMessage, error) {
	filter := bson.M{"email": email}

	cursor, err := dbRepo.AppConfig.MongoDatabase.Collection("imageData").Find(context.Background(), filter)
	if err != nil {
		log.Println("email not found")
		return nil, err
	}
	var jsonArray []json.RawMessage
	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}

		jsonBytes, err := json.Marshal(document)
		if err != nil {
			log.Fatal(err)
		}

		jsonArray = append(jsonArray, json.RawMessage(jsonBytes))
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return jsonArray, nil

}
