package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)

	authToken := os.Getenv("SLACK_AUTH_TOKEN")
	if authToken == "" {
		logger.Fatal("No value found for environment variable: \"SLACK_AUTH_TOKEN\"")
	}

	slackAPI := slack.New(authToken)
	slack.SetLogger(logger)
	slackAPI.SetDebug(false)

	postChannel := os.Getenv("SLACK_CHANNEL")
	// If no destination channel has been defined, use slackbot
	if postChannel == "" {
		IMChannels, err := slackAPI.GetIMChannels()
		if err != nil {
			logger.Printf("%s\n", err)
		}

		for _, IMChannel := range IMChannels {
			if IMChannel.User == "USLACKBOT" {
				postChannel = IMChannel.ID
				break
			}
		}
	}

	sleepTimeInt := 300 // Default sleep time as 5 minutes
	sleepTimeString := os.Getenv("SLACK_CHECK_INTERVAL")
	if sleepTimeString != "" {
		var err error
		sleepTimeInt, err = strconv.Atoi(sleepTimeString)
		if err != nil {
			logger.Fatal(err)
		}
	}
	sleepTime := time.Duration(sleepTimeInt) * time.Second

	currentChannels := []string{}
	for {
		previousChannels := currentChannels
		currentChannels = []string{}
		channels, err := slackAPI.GetChannels(true)
		if err != nil {
			logger.Fatal(err)
		}
		for _, channel := range channels {
			if len(previousChannels) != 0 && !stringInSlice(channel.ID, previousChannels) {
				messageParams := slack.PostMessageParameters{}
				messageParams.Username = "New Slack Channel"
				messageParams.IconEmoji = ":bell:"
				messageText := fmt.Sprintf("Channel #%s has just become available", channel.Name)
				_, _, err := slackAPI.PostMessage(postChannel, messageText, messageParams)
				if err != nil {
					logger.Printf("%s\n", err)
				}
				logger.Println(messageText)
			}
			currentChannels = append(currentChannels, channel.ID)
		}
		time.Sleep(sleepTime)
	}
}
