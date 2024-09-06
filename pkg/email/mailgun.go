package email

import "github.com/mailgun/mailgun-go"

type (
	MailgunTransport struct {
		Domain string `yaml:"domain"`
		APIKey string `yaml:"api_key"`
	}

	MailgunClient struct {
		mailgun.Mailgun
		opts *ClientOptions[MailgunTransport]
	}
)

func NewMailgunClient(opts *ClientOptions[MailgunTransport]) *MailgunClient {
	return &MailgunClient{
		Mailgun: mailgun.NewMailgun(opts.Transport.Domain, opts.Transport.APIKey),
		opts:    opts,
	}
}
