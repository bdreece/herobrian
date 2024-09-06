package email

import "context"

type (
	Address struct {
		Name  string
		Email string
	}

	Message struct {
		To      Address
		Subject string
		Text    string
		HTML    string
	}

	ClientOptions[T any] struct {
		From      Address
		Transport T
	}

	Client interface {
		Send(context.Context, *Message) error
	}
)
