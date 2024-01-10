package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

const (
	Markdown = tgbotapi.ModeMarkdown
	Html     = tgbotapi.ModeHTML
)

var parseMode = Html
var name = "telegram"

type Telegram struct {
	client  *tgbotapi.BotAPI
	chatIDs []int64
}

// For more information about telegram api token:
// https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api#NewBotAPI
func New(apiToken string, serviceName string) (*Telegram, error) {
	if serviceName != "" {
		name = serviceName
	}
	client, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}

	t := &Telegram{
		client:  client,
		chatIDs: []int64{},
	}

	return t, nil
}

// For example allowing you to use NewBotAPIWithClient:
// https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api#NewBotAPIWithClient
func (t *Telegram) SetClient(client *tgbotapi.BotAPI) *Telegram {
	t.client = client
	return t
}

// SetParseMode sets the parse mode for the message body.
// https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api#pkg-constants
func (t *Telegram) SetParseMode(mode string) *Telegram {
	parseMode = mode
	return t
}

func (t *Telegram) AddReceivers(chatIDs ...int64) *Telegram {
	t.chatIDs = append(t.chatIDs, chatIDs...)
	return t
}

func GetName() string {
	return name
}

// object position: message, parse mode
func (t Telegram) Send(ctx context.Context, object ...interface{}) error {
	message := fmt.Sprintf("%v", object[0])

	msg := tgbotapi.NewMessage(0, message)
	msg.ParseMode = parseMode

	for _, chatID := range t.chatIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg.ChatID = chatID
			_, err := t.client.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Telegram chat '%d'", chatID)
			}
		}
	}

	return nil
}
