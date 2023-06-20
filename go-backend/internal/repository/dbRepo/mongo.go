package dbrepo

import (
	"context"
	"github/toothsy/go-background-job/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (dbRepo *mongoDBRepo) SignUp(user *models.UserPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := dbRepo.AppConfig.MongoDatabase.Collection("userData").InsertOne(ctx, user)
	if err != nil {
		log.Println("\n\n could not insert, err is ", err)
	}

}

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
	log.Println("Fetched User ", foundUser)
	return foundUser, nil
}

func (dbRepo *mongoDBRepo) UpdateUserVerification(user *models.UserPayload) {
	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": bson.M{"isverified": true}}
	result, err := dbRepo.AppConfig.MongoDatabase.Collection("userData").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Did not update User, as user not found")
	}
	log.Println("Updated User :", result.MatchedCount)

}
