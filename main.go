package main

import (
	"felix-notification/example"
	"felix-notification/notification"
)

func main() {
	// main code here
	notification := notification.Notification{}
	notification.AddService(example.MailExample())
}
