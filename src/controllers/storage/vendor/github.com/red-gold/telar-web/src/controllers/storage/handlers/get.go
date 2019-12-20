package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	handler "github.com/openfaas-incubator/go-function-sdk"
	"github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	appConfig "github.com/red-gold/telar-web/src/controllers/storage/config"
	uuid "github.com/satori/go.uuid"
)

var redisClient *redis.Client

func init() {

}

// GetFileHandle a function invocation
func GetFileHandle(db interface{}) func(http.ResponseWriter, *http.Request, server.Request) (handler.Response, error) {

	return func(w http.ResponseWriter, r *http.Request, req server.Request) (handler.Response, error) {

		storageConfig := &appConfig.StorageConfig

		// Initialize Redis Connection
		if redisClient == nil {

			redisPassword, redisErr := utils.ReadSecret("redis-pwd")

			if redisErr != nil {
				fmt.Printf("\n\ncouldn't get payload-secret: %s\n\n", redisErr.Error())
			}
			fmt.Println(storageConfig.RedisAddress)
			fmt.Println(redisPassword)
			redisClient = redis.NewClient(&redis.Options{
				Addr:     storageConfig.RedisAddress,
				Password: redisPassword,
				DB:       0,
			})
			pong, err := redisClient.Ping().Result()
			fmt.Println(pong, err)
		}

		fmt.Println("File Upload Endpoint Hit")
		// params from /storage/:uid/:dir/:name
		dirName := req.GetParamByName("dir")
		if dirName == "" {
			errorMessage := fmt.Sprintf("Directory name is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("dirNameRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Directory name: %s", dirName)

		// params from /storage/:dir
		fileName := req.GetParamByName("name")
		if fileName == "" {
			errorMessage := fmt.Sprintf("File name is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("fileNameRequired", errorMessage)}, nil
		}
		fmt.Printf("\n File name: %s", fileName)

		userId := req.GetParamByName("uid")
		if userId == "" {
			errorMessage := fmt.Sprintf("User Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("userIdRequired", errorMessage)}, nil
		}

		fmt.Printf("\n User ID: %s", userId)
		userUUID, uuidErr := uuid.FromString(userId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}

		objectName := fmt.Sprintf("%s/%s/%s", userUUID, dirName, fileName)
		objectExpireKey := fmt.Sprintf("expire:%s", objectName)

		objectURLKey := fmt.Sprintf("url:%s", objectName)
		expireTime := time.Now().Add(7100 * time.Second)
		expireUnix := utils.TimeUnix(expireTime)
		expireDate, expireKeyErr := redisClient.Get(objectExpireKey).Result()
		if expireKeyErr != nil {
			errorMessage := fmt.Sprintf("Get expire by key Error %s", expireKeyErr.Error())
			fmt.Println(errorMessage)
			// return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("expireKeyError", errorMessage)}, nil
		}

		isExpired := true
		downloadURL := ""
		if expireDate != "" {
			if parsedUnix, err := strconv.ParseInt(expireDate, 10, 64); err == nil {
				fmt.Println(parsedUnix)
				isExpired = utils.IsTimeExpired(parsedUnix, 7100)
			} else {
				errorMessage := fmt.Sprintf("Expire time is not integer %s : %s", expireDate, err.Error())
				return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("parseExpireUnixError", errorMessage)}, nil
			}

		}

		if !isExpired {
			var urlKeyErr error
			downloadURL, urlKeyErr = redisClient.Get(objectURLKey).Result()
			if urlKeyErr != nil {
				errorMessage := fmt.Sprintf("Get URL by key Error %s", urlKeyErr.Error())
				fmt.Println(errorMessage)
				// return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("urlKeyError", errorMessage)}, nil
			}
		} else {
			// Generate download URL
			var urlErr error
			downloadURL, urlErr = generateV4GetObjectSignedURL(storageConfig.BucketName, objectName, storageConfig.StorageSecretPath)
			if urlErr != nil {
				fmt.Println(urlErr.Error())
			}

			redisErr := redisClient.MSet(objectURLKey, downloadURL, objectExpireKey, expireUnix).Err()
			if redisErr != nil {
				errorMessage := fmt.Sprintf("Redis Error %s", redisErr.Error())
				return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("redisError", errorMessage)}, nil
			}
		}

		code := 302 // Permanent redirect, request with GET method
		if r.Method != http.MethodGet {
			// Temporary redirect, request with same method
			// As of Go 1.3, Go does not support status code 308.
			code = 307
		}

		http.Redirect(w, r, downloadURL, code)

		return handler.Response{
			StatusCode: http.StatusTemporaryRedirect,
		}, nil

	}

}
