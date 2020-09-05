package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// TODO: just here for test purpose.
var slackToken = os.Getenv("SLACK_TOKEN")
var slackVerifToken = os.Getenv("SLACK_VERIF_TOKEN")
var messageResponse = `Hei :wave: This is an automatic message. 

Unfortunaly, I don't accept private messages.

If you have a question regarding your studies: bocal@hive.fi, if you want to contact me: titus@hive.fi.

Thank you.
`

var api = slack.New(slackToken)

func main() {
	http.HandleFunc("/events", EventHandler)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}
func EventHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: slackVerifToken},
		),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if ev.SubType != "" || isToMyself(ev.Channel) || isAdminOrBot(ev.User) {
				return
			}

			log.Println("[INFO] new message event.")

			api.PostMessage(
				ev.Channel,
				slack.MsgOptionText(messageResponse, false),
			)
			log.Println("[INFO] handled.")
		}
	}
}

func isAdminOrBot(userID string) bool {
	user, err := api.GetUserInfo(userID)
	if err != nil {
		log.Printf("[ERROR] Can't get user info: %v\n", err)
		return true
	}
	return (user.IsAdmin || user.IsBot)
}

func isToMyself(channelID string) bool {
	users, _, err := api.GetUsersInConversation(
		&slack.GetUsersInConversationParameters{ChannelID: channelID},
	)
	if err != nil {
		log.Printf("[ERROR] Can't get conversation users lists: %v\n", err)
	}
	return len(users) == 1
}
