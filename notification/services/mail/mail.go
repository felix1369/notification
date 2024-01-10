package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

type Mail struct {
	usePlainText      bool
	senderAddress     string
	smtpHostAddr      string
	smtpAuth          smtp.Auth
	receiverAddresses []string
	ccAddresses       []string
	bccAddresses      []string
}

var name = "mail"

func New(senderAddress, smtpHostAddress string, serviceName string) *Mail {
	if serviceName != "" {
		name = serviceName
	}
	return &Mail{
		usePlainText:      false,
		senderAddress:     senderAddress,
		smtpHostAddr:      smtpHostAddress,
		receiverAddresses: []string{},
	}
}

type BodyType int

const (
	PlainText BodyType = iota
	HTML
)

// For more information about smtp authentication, see here:
//	-> https://pkg.go.dev/net/smtp#PlainAuth
// Example values: "", "test@gmail.com", "password123", "smtp.gmail.com"
func (m *Mail) AuthenticateSMTP(identity, userName, password string) {
	m.smtpAuth = smtp.PlainAuth(identity, userName, password, strings.Split(m.smtpHostAddr, ":")[0])
}

func (m *Mail) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

func (m *Mail) AddCarbonCopy(addresses ...string) {
	m.ccAddresses = append(m.ccAddresses, addresses...)
}

func (m *Mail) AddBlindCarbonCopy(addresses ...string) {
	m.bccAddresses = append(m.bccAddresses, addresses...)
}

// BodyFormat can be used to specify the format of the body.
// Default BodyType is HTML.
func (m *Mail) BodyFormat(format BodyType) {
	switch format {
	case PlainText:
		m.usePlainText = true
	default:
		m.usePlainText = false
	}
}

func (m *Mail) newEmail(subject, message string) *email.Email {
	msg := &email.Email{
		To:      m.receiverAddresses,
		Cc:      m.ccAddresses,
		Bcc:     m.bccAddresses,
		From:    m.senderAddress,
		Subject: subject,
		Headers: textproto.MIMEHeader{},
	}

	if m.usePlainText {
		msg.Text = []byte(message)
	} else {
		msg.HTML = []byte(message)
	}
	return msg
}

func (m Mail) GetName() string {
	return name
}

// object position: subject, body
func (m Mail) Send(ctx context.Context, object ...interface{}) error {
	subject := fmt.Sprintf("%v", object[0])
	message := fmt.Sprintf("%v", object[0])
	msg := m.newEmail(subject, message)

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		err = msg.Send(m.smtpHostAddr, m.smtpAuth)
		if err != nil {
			err = errors.Wrap(err, "failed to send mail")
		}
	}

	return err
}
