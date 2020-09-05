package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jgengo/AutoResponder/internal/slacker"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// TODO: remove only used for test
var slackVerifToken = os.Getenv("SLACK_VERIF_TOKEN")
var Api *slack.Client

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
			if ev.SubType != "" || slacker.IsToMyself(ev.Channel) || slacker.IsAdminOrBot(ev.User) {
				return
			}

			log.Println("[INFO] new message event.")

			slacker.PostMessage(ev.Channel)

			log.Println("[INFO] handled.")
		}
	}
}
