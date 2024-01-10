package notification

import "context"

type NotificationRequest struct {
	Context context.Context
	Object  interface{}
}
