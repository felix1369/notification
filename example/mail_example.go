package example

import (
	"context"
	"felix-notification/notification/services/mail"
	"fmt"
)

func MailExample() *mail.Mail {
	mail := mail.New("from@example.com", "sandbox.smtp.mailtrap.io:25", "mail1")
	mail.AuthenticateSMTP("", "80b2c2b3fda7f8", "3dfe28a193bc59")
	mail.AddReceivers("to@example.com")
	return mail
}
func MailSend(mail *mail.Mail) {
	pesan := `<html>
				<head>
				<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3DUTF-8">
				</head>
				<body style=3D"font-family: sans-serif;">
				<div style=3D"display: block; margin: auto; max-width: 600px;" class=3D"main">
					<h1 style=3D"font-size: 18px; font-weight: bold; margin-top: 20px">Congrats for sending test email with Mailtrap!</h1>
					<p>If you are viewing this email in your inbox =E2=80=93 the integration works.</p>
					<img alt=3D"Inspect with Tabs" src=3D"https://assets-examples.mailtrap.io/integration-examples/welcome.png" style=3D"width: 100%;">
					<p>Now send your email using our SMTP server and integration of your choice!</p>
					<p>Good luck! Hope it works.</p>
				</div>
				<!-- Example of invalid for email html/css, will be detected by Mailtrap: -->
				<style>
					.main { background-color: white; }
					a:hover { border-left-width: 1em; min-height: 2em; }
				</style>
				</body>
			</html>`
	err := mail.Send(context.Background(), pesan)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("msg sent")
}
