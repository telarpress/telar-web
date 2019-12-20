package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alexellis/hmac"
)

func main() {
	payloadSecret := "70bc03d019a6250e5731655e8d2528c68dd27318"

	bytesReq, _ := json.Marshal(struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}{
		Type: "SET_MESSAGE",
		Payload: struct {
			Message string `json:"message"`
		}{
			Message: "test message beauty",
		},
	})
	fmt.Println(string(bytesReq))

	// bytesReq := []byte(`{"type": "SET_MESSAGE","payload": {"message": "test message"}}`)
	digest := hmac.Sign(bytesReq, []byte(payloadSecret))
	fmt.Println("sha1=" + hex.EncodeToString(digest))
	bufferReader := bytes.NewBuffer(bytesReq)

	request, _ := http.NewRequest(http.MethodPost, "http://localhost:3000/dispatch/ef54ada2-3e8f-4c88-8406-7eeaf4d9fc40", bufferReader)
	request.Header.Set("Content-type", "application/json")
	request.Header.Add("Http_Hmac", "sha1="+hex.EncodeToString(digest))
	c := http.Client{
		Timeout: time.Second * 3,
	}
	response, err := c.Do(request)

	if err == nil {
		if response.Body != nil {
			defer response.Body.Close()
			fmt.Println("BODY CONTENT")

			bodyBytes, bErr := ioutil.ReadAll(response.Body)
			if bErr != nil {
				log.Fatal(bErr)
			}
			log.Println(string(bodyBytes))
		}
	} else {
		fmt.Println(err.Error())
	}
}
