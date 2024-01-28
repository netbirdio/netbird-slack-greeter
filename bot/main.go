package main

import (
	"bytes"
	sl "log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"

	log "github.com/sirupsen/logrus"
)

const (
	bugsIssuesEtcTemplate = "Hello, <@{{.User}}>, and welcome to the <#{{.Channel}}> channel of NetBird Slack!☁️❗\n \n" +
		"This channel is dedicated to queries on configuring your NetBird account, issues running NetBird clients, " +
		"possible feature requests, and sharing use cases.\n\nFor prompt and effective assistance with issues, " +
		"*please provide detailed information as outlined in our bug/issue reporting template*: " +
		"https://docs.netbird.io/how-to/report-bug-issues#reporting-template\n\n*Note*: _We prioritize support for queries " +
		"posted in the appropriate channel and those that follow our template guidelines._\n🌟 Your contributions are " +
		"essential to us! Feel free to report bugs or suggest features via GitHub issues: " +
		"https://github.com/netbirdio/netbird/issues"
	selfHostedTemplate = "Hello, <@{{.User}}>, and welcome to the <#{{.Channel}}> channel of " +
		"NetBird Slack! 🛠️\n\nThis channel is dedicated to discussions, support, and sharing issues specifically related " +
		"to our self-hosted deployments.\n\nFor prompt and effective assistance, *please provide detailed information as outlined " +
		"in our bug/issue reporting template*: https://docs.netbird.io/how-to/report-bug-issues#reporting-template\n\n*Note*: " +
		"_We prioritize support for queries posted in the appropriate channel and those that follow our template guidelines._\n\n🌟 " +
		"Your contributions are essential to us! Feel free to report bugs or suggest features via GitHub issues: " +
		"https://github.com/netbirdio/netbird/issues"
	newUserTemplate = "Hello, <@{{.User}}>,  and welcome to NetBird's community Slack!🌐\n\n" +
		"You will find updates and notifications from the NetBird team in our <#C028VPB34NB> " +
		"channel.\n\n🔍 Are you encountering issues or need support?\n • " +
		"For queries on how to configure your network or issues running netbird clients, " +
		"join <#C02KHAE8VLZ>\n\n  • For self-hosted issues, like setting up services or client communication, " +
		"join <#C05T5K65X7U>\n\nFor prompt and effective assistance, *please provide detailed information as " +
		"outlined in our bug/issue reporting template*: https://docs.netbird.io/how-to/report-bug-issues#reporting-template"
)

type templateInput struct {
	User    string
	Channel string
}

func initLog() *log.Logger {
	logFormatter := new(log.JSONFormatter)
	logFormatter.TimestampFormat = time.RFC3339 // or RFC3339
	logFormatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
		return "", fileName
	}
	log.SetFormatter(logFormatter)
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	logger := log.New()
	logger.SetFormatter(logFormatter)
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.TraceLevel)
	return logger
}

func main() {
	logger := initLog()
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		panic("SLACK_APP_TOKEN must be set.\n")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		panic("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		panic("SLACK_BOT_TOKEN must be set.\n")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		panic("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(false),
		slack.OptionLog(sl.New(logger.Writer(), "api: ", sl.Lshortfile|sl.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(sl.New(logger.Writer(), "socketmode: ", sl.Lshortfile|sl.LstdFlags)),
	)

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)

	socketmodeHandler.Handle(socketmode.EventTypeEventsAPI, getMiddlewareEventsAPI(api))

	socketmodeHandler.RunEventLoop()
}

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	log.Info("connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	log.Error("connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	log.Info("connected to Slack with Socket Mode.")
}

func getMiddlewareEventsAPI(api *slack.Client) func(evt *socketmode.Event, client *socketmode.Client) {
	return func(evt *socketmode.Event, client *socketmode.Client) {
		eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
		if !ok {
			log.Tracef("event ignored %+v", evt)
			return
		}

		log.Infof("event received: %+v\n", eventsAPIEvent)

		client.Ack(*evt.Request)

		switch eventsAPIEvent.Type {
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MemberJoinedChannelEvent:
				log.Debugf("user %q joined to channel %q", ev.User, ev.Channel)
				templateText := getChannelTemplate(ev.Channel)
				err := postMessage(api, true, ev.User, ev.Channel, templateText)
				if err != nil {
					log.Errorf("got an error when posting a message: %s\n", err)
				}
			case *slackevents.TeamJoinEvent:
				log.Debugf("user %q joined to slack, its name is: %q", ev.User.ID, ev.User.Name)
				channel, _, _, err := api.OpenConversation(&slack.OpenConversationParameters{Users: []string{ev.User.ID}})
				if err != nil {
					log.Errorf("got an error when opening a conversation: %s\n", err)
					return
				}
				err = postMessage(api, false, ev.User.ID, channel.ID, newUserTemplate)
				if err != nil {
					log.Errorf("got an error when posting a message: %s\n", err)
				}
			}
		default:
		}
	}
}

func getChannelTemplate(channel string) string {
	switch channel {
	case "C02KHAE8VLZ":
		return bugsIssuesEtcTemplate
	case "C05T5K65X7U":
		return selfHostedTemplate
	default:
		return ""
	}
}

func postMessage(api *slack.Client, isChannelMSG bool, user, channel, templateText string) error {
	if templateText == "" {
		log.Infof("no template for channel %s", channel)
		return nil
	}

	text, err := parseText(templateInput{
		User:    user,
		Channel: channel,
	}, templateText)
	if err != nil {
		return err
	}

	options := []slack.MsgOption{slack.MsgOptionText(text, false), slack.MsgOptionAsUser(true)}
	if isChannelMSG {
		options = append([]slack.MsgOption{slack.MsgOptionPostEphemeral(user)}, options...)
	}

	channelID, timestamp, err := api.PostMessage(channel, options...)
	if err != nil {
		return err
	}
	log.Infof("message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

func parseText(input templateInput, textTemplate string) (string, error) {
	tmpl, err := template.New("msgs").Parse(textTemplate)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, input)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
