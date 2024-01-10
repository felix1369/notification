package notification

import "context"

type NotificationService interface {
	GetName() string
	Send(ctx context.Context, object ...interface{}) error // object position: message, receiver, reserved param 1, reserved param 2 etc
}
