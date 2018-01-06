# slack-channelwatcher

## About

This is a simple Go application which will actively monitor your Slack channels,
and send you a notification in the event that a new channel is created.

It's currently written in such a way that it would use your user API token in
Slack to monitor the channels which are available to you, and notify you when a
new one is created/un-archived.

## Configuration

Configuration is achieved through environment variables, there are a few items
which are configurable:

- **SLACK_AUTH_TOKEN** - Your Slack API authentication token. For more
  information on creating this, [check Slack's documentation](https://get.slack.help/hc/en-us/articles/215770388-Create-and-regenerate-API-tokens).
  (**Required**)
- **SLACK_CHECK_INTERVAL** *(seconds)* - Define the frequency at which the app
  will poll for channel changes. (**Default:** 300)
- **SLACK_CHANNEL** - Destination channel for notifications. It's highly
  recommended to use the channel identifier rather than the channel name, in
  case of the event that the channel name should change. (**Default:**
  `@slackbot`)

## Usage

Install and run on your local host (this assumes the Go `bin` directory is in your PATH):
```
go get github.com/nathandines/slack-channelwatcher
SLACK_AUTH_TOKEN='<slack_token_here>' slack-channelwatcher
```

Run in Docker:
```
docker run -e 'SLACK_AUTH_TOKEN=<slack_token_here>' nathandines/slack-channelwatcher:latest
```
