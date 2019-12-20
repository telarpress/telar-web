package handlers

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	cf "github.com/red-gold/telar-web/src/controllers/users/auth/config"
)

type HomepageTokens struct {
	AccessToken string
	Login       string
}

// HomepageHandler shows the homepage
func HomepageHandler(w http.ResponseWriter, r *http.Request, req server.Request) (handler.Response, error) {
	config := cf.AuthConfig
	keydata, err := ioutil.ReadFile("config.PublicKeyPath")
	if err != nil {
		log.Fatalf("unable to read path: %s, error: %s", "config.PublicKeyPath", err.Error())
	}

	publicKey, keyErr := jwt.ParseECPublicKeyFromPEM(keydata)
	if keyErr != nil {
		log.Fatalf("unable to parse public key: %s", keyErr.Error())
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Println("No cookie found.")
		prettyURL := utils.GetPrettyURLf(config.BaseRoute)

		http.Redirect(w, r, prettyURL+"/login/?r="+r.URL.Path, http.StatusTemporaryRedirect)
		return handler.Response{}, nil
	}

	parsed, parseErr := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if parseErr != nil {
		log.Println(parseErr, cookie.Value)
		return handler.Response{
			Body: []byte("Unable to decode cookie, please clear your cookies and sign-in again"),
		}, nil
	}
	log.Printf("Parsed JWT: %v", parsed)

	tmpl, err := template.ParseFiles("./html_template/home.html")

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, HomepageTokens{
		AccessToken: "Unavailable",
		Login:       "Unknown",
	})

	if err != nil {
		log.Panic("Error executing template: ", err)
	}

	return handler.Response{
		Body: tpl.Bytes(),
	}, nil
}
