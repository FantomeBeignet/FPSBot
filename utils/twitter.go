package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	twitter "github.com/g8rswimmer/go-twitter/v2"
	"github.com/joho/godotenv"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func GetTwtProfilePicture(profile string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	u, err := url.Parse(profile)
	if err != nil {
		log.Fatal(err)
	}
	name := u.Path[1:]
	client := &twitter.Client{
		Authorizer: authorize{
			Token: os.Getenv("TWITTER_TOKEN"),
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{
			twitter.UserFieldProfileImageURL,
		},
	}
	userResponse, err := client.UserNameLookup(context.Background(), []string{name}, opts)
	if err != nil {
		log.Fatal(err)
	}
	return userResponse.Raw.Users[0].ProfileImageURL
}
