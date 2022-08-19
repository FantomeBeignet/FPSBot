package utils

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetYtProfilePicture(profile string) string {
	u, err := url.Parse(profile)
	if err != nil {
		log.Fatal(err)
	}
	var channelId string
	linkType := strings.Split(u.Path, "/")[1]
	if linkType == "channel" {
		channelId = strings.Split(u.Path, "/")[2]
	} else {
		username := strings.Split(u.Path, "/")[2]
		channelId = GetChannelId(username)
	}
	return getChannelProfilePicture(channelId)
}

func GetChannelId(profile string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(googleAPIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
	call := service.Channels.List([]string{"contentDetails"})
	call = call.ForUsername(profile)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
	}
	return response.Items[0].Id
}

func getChannelProfilePicture(channelId string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(googleAPIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
	call := service.Channels.List([]string{"snippet"}).Id(channelId)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
	}
	return response.Items[0].Snippet.Thumbnails.Default.Url
}
