package domain

type Notification struct {
	Code    string
	Payload interface{}
}

type NotificationService interface {
	Publish(notification Notification)
	Subscribe(handler func(notification Notification))
}
