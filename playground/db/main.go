package main

import (
	"context"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

const uri = "mongodb://uyca00tvvj6liagozmao:PhLr52XMWONFDDcbEH26@bbnjlclnm5piolw-mongodb.services.clever-cloud.com:27017/bbnjlclnm5piolw"

// const uri = "mongodb+srv://telar_user:pass@cluster0-l6ojz.mongodb.net/test?retryWrites=true&w=majority"

type UserProfile struct {
	ObjectId       uuid.UUID `json:"objectId" bson:"objectId"`
	FullName       string    `json:"fullName" bson:"fullName"`
	Avatar         string    `json:"avatar" bson:"avatar"`
	Banner         string    `json:"banner" bson:"banner"`
	TagLine        string    `json:"tagLine" bson:"tagLine"`
	CreatedDate    int64     `json:"created_date" bson:"created_date"`
	LastUpdated    int64     `json:"last_updated" bson:"last_updated"`
	Email          string    `json:"email" bson:"email"`
	Birthday       int64     `json:"birthday" bson:"birthday"`
	WebUrl         string    `json:"webUrl" bson:"webUrl"`
	CompanyName    string    `json:"companyName" bson:"companyName"`
	VoteCount      int64     `json:"voteCount" bson:"voteCount"`
	ShareCount     int64     `json:"shareCount" bson:"shareCount"`
	FollowCount    int64     `json:"followCount" bson:"followCount"`
	FollowerCount  int64     `json:"followerCount" bson:"followerCount"`
	PostCount      int64     `json:"postCount" bson:"postCount"`
	FacebookId     string    `json:"facebookId" bson:"facebookId"`
	InstagramId    string    `json:"instagramId" bson:"instagramId"`
	TwitterId      string    `json:"twitterId" bson:"twitterId"`
	AccessUserList []string  `json:"accessUserList" bson:"accessUserList"`
}

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Print(err.Error())
	}

	userUUID, _ := uuid.FromString("ef54ada2-3e8f-4c88-8406-7eeaf4d9fc40")

	// filter := struct {
	// 	ObjectId uuid.UUID `bson:"objectId"`
	// }{
	// 	ObjectId: userUUID,
	// }
	filter := bson.M{"objectId": userUUID}
	client.Connect(ctx)
	// _ = client.Ping(ctx, readpref.Primary())
	collection := client.Database("bbnjlclnm5piolw").Collection("userProfile")

	var result UserProfile
	for index := 0; index < 3; index++ {
		start := time.Now()
		err = collection.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		duration := time.Since(start)
		fmt.Printf("\nProfile %v\n", result)
		fmt.Println("***********************************")
		fmt.Printf("\n Duration: %f\n", duration.Seconds())
		fmt.Println("***********************************")
		delaySecond(5)

	}

}
