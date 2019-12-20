package handlers

import (
	"fmt"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

func generateV4GetObjectSignedURL(bucketName string, objectName string, serviceAccount string) (string, error) {
	// [START storage_generate_signed_url_v4]
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(120 * time.Hour), // 5 days , 7200 seconds
	}

	u, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("Unable to generate a signed URL: %v", err)
	}

	fmt.Println("Generated GET signed URL:")
	fmt.Printf("%q\n", u)
	fmt.Println("You can use this URL with any user agent, for example:")
	fmt.Printf("curl %q\n", u)
	// [END storage_generate_signed_url_v4]
	return u, nil
}
