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
	logger := log.New(os.Stdout, "slack-bot: ", log.LstdFlags)

	authToken := os.Getenv("SLACK_AUTH_TOKEN")
	autoJoinChannelVar := os.Getenv("SLACK_AUTO_JOIN_CHANNEL")
	sleepTimeString := os.Getenv("SLACK_CHECK_INTERVAL")
	postChannel := os.Getenv("SLACK_CHANNEL")

	if authToken == "" {
		logger.Fatal("No value found for environment variable: \"SLACK_AUTH_TOKEN\"")
	}

	autoJoinChannel, err := strconv.ParseBool(autoJoinChannelVar)
	if err != nil {
		if autoJoinChannelVar == "" {
			autoJoinChannel = false
		} else {
			logger.Fatal("\"SLACK_AUTO_JOIN_CHANNEL\" must be a boolean value")
		}
	}

	sleepTimeInt, err := strconv.Atoi(sleepTimeString)
	if err != nil {
		if sleepTimeString == "" {
			sleepTimeInt = 300 // Default sleep time as 5 minutes
		} else {
			logger.Fatal(err)
		}
	}

	sleepTime := time.Duration(sleepTimeInt) * time.Second

	slackAPI := slack.New(authToken)
	slack.SetLogger(logger)
	slackAPI.SetDebug(false)

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

	logger.Println("Monitoring Slack channels...")

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
				messageParams.Username = "Slack Channel Watcher"
				messageParams.IconEmoji = ":eyes:"
				messageText := fmt.Sprintf("New channel: <#%s|%s>", channel.ID, channel.Name)
				_, _, err := slackAPI.PostMessage(postChannel, messageText, messageParams)
				if err != nil {
					logger.Printf("%s\n", err)
				}
				logger.Println(messageText)
				if autoJoinChannel {
					slackAPI.JoinChannel(channel.Name)
					logger.Printf("Joined channel: #%s", channel.Name)
				}
			}
			currentChannels = append(currentChannels, channel.ID)
		}
		time.Sleep(sleepTime)
	}
}
