package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/go-redis/redis/v7"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-19712.c12.us-east-1-4.ec2.cloud.redislabs.com:19712",
		Password: "zBlSkiMsP61RWOwHa0xla7sH628TXkF5", // no password set
		DB:       0,                                  // use default DB
	})

	pong, err := redisClient.Ping().Result()

	fmt.Println(pong, err)

	objectNameRedis := "userId/dir/picId"
	objectExpireKey := fmt.Sprintf("expire:%s", objectNameRedis)

	URL := "http://storage.com/userId/dir/picId"

	objectURLKey := fmt.Sprintf("url:%s", objectNameRedis)
	expireTime := time.Now().Add(5 * time.Second)

	err = redisClient.MSet(objectURLKey, URL, objectExpireKey, expireTime).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get(objectExpireKey).Result()
	if err != nil {
		panic(err)
	}

	t, err := time.Parse(time.RFC3339, val)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	remainder := t.Sub(time.Now())
	fmt.Printf("\nremainder: %v calc : %v", remainder, (remainder.Seconds() + 10))
	fmt.Printf("\nIs expired: %t\n", !((remainder.Seconds() + 10) > 0))
	fmt.Println(objectExpireKey, val)

	x := "hello"
	if x == "hello" {
		return
	}
	config := &firebase.Config{
		StorageBucket: "resume-web-app.appspot.com",
	}
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}

	// [START upload_file]
	f, err := os.Open("notes.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	wc := bucket.Object("directory/test-file.txt").NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		fmt.Println(err.Error())
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err.Error())
	}

}
