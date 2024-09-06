package email

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/config"
)

type (
	SendGridOptions struct {
		APIKey string `yaml:"api_key"`
		From   struct {
			Name  string `yaml:"name"`
			Email string `yaml:"email"`
		} `yaml:"from"`
	}

	SendGridClient struct {
		from *mail.Email
		*sendgrid.Client
	}
)

func (c *SendGridClient) Send(ctx context.Context, msg *Message) error {
	to := mail.NewEmail(msg.To.Name, msg.To.Email)
	mailMsg := mail.NewSingleEmail(c.from, msg.Subject, to, msg.Text, msg.HTML)
	res, err := c.Client.SendWithContext(ctx, mailMsg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("received invalid status code: %d", res.StatusCode)
	}

	return nil
}

func ConfigureSendGrid(provider config.Provider) (*SendGridOptions, error) {
	opts := new(SendGridOptions)
	if err := provider.Get("sendgrid").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure sendgrid options: %w", err)
	}

	return opts, nil
}

func NewSendGridClient(opts *SendGridOptions) *SendGridClient {
	return &SendGridClient{
		Client: sendgrid.NewSendClient(opts.APIKey),
		from:   mail.NewEmail(opts.From.Name, opts.From.Email),
	}
}
