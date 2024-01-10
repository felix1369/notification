package slack

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type slackClient interface {
	PostMessageContext(ctx context.Context, channelID string, options ...slack.MsgOption) (string, string, error)
}

var _ slackClient = new(slack.Client)
var name = "slack"

type Slack struct {
	client     slackClient
	channelIDs []string
}

// For more information about slack api token: https://pkg.go.dev/github.com/slack-go/slack#New
func New(apiToken string, serviceName string) *Slack {
	if serviceName != "" {
		name = serviceName
	}
	client := slack.New(apiToken)

	s := &Slack{
		client:     client,
		channelIDs: []string{},
	}

	return s
}

func (s *Slack) AddReceivers(channelIDs ...string) {
	s.channelIDs = append(s.channelIDs, channelIDs...)
}

func GetName() string {
	return name
}

// need a slack app with the chat:write.public and chat:write permissions. see https://api.slack.com/
// object position: message
func (s Slack) Send(ctx context.Context, object ...interface{}) error {
	message := fmt.Sprintf("%v", object[0])

	for _, channelID := range s.channelIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			id, timestamp, err := s.client.PostMessageContext(
				ctx,
				channelID,
				slack.MsgOptionText(message, false),
			)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Slack channel '%s' at time '%s'", id, timestamp)
			}
		}
	}

	return nil
}
