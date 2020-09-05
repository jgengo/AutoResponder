package slacker

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// TODO: just here for test purpose.
var slackToken = os.Getenv("SLACK_TOKEN")
var messageResponse = `Hei :wave: This is an automatic message. 

Unfortunaly, I don't accept private messages.

If you have a question regarding your studies: bocal@hive.fi, if you want to contact me: titus@hive.fi.

Thank you.
`

var api = slack.New(slackToken)

func IsAdminOrBot(userID string) bool {
	user, err := api.GetUserInfo(userID)
	if err != nil {
		log.Printf("[ERROR] Can't get user info: %v\n", err)
		return true
	}
	return (user.IsAdmin || user.IsBot)
}

func IsToMyself(channelID string) bool {
	users, _, err := api.GetUsersInConversation(
		&slack.GetUsersInConversationParameters{ChannelID: channelID},
	)
	if err != nil {
		log.Printf("[ERROR] Can't get conversation users lists: %v\n", err)
	}
	return len(users) == 1
}

func PostMessage(channelID string) {
	api.PostMessage(
		channelID,
		slack.MsgOptionText(messageResponse, false),
	)
}
