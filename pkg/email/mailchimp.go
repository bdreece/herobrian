package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/config"
	"go.uber.org/multierr"
)

type (
	MailchimpTransport struct {
		APIKey string `yaml:"api_key"`
	}

	MailchimpClient struct {
		opts *ClientOptions[MailchimpTransport]

		http.Client
	}
	mailchimpRequest struct {
		Key     string           `json:"key"`
		Message mailchimpMessage `json:"message"`
	}
	mailchimpResponse struct {
		Email        string `json:"email"`
		Status       string `json:"status"`
		RejectReason string `json:"reject_reason"`
	}
	mailchimpMessage struct {
		HTML      string               `json:"html"`
		Text      string               `json:"text"`
		Subject   string               `json:"subject"`
		FromName  string               `json:"from_name"`
		FromEmail string               `json:"from_email"`
		To        []mailchimpRecipient `json:"to"`
	}
	mailchimpRecipient struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func (client *MailchimpClient) Send(ctx context.Context, msg *Message) (err error) {
	const uri string = "https://mandrillapp.com/api/1.0/messages/send"
	var buf bytes.Buffer

	reqBody := mailchimpRequest{
		Key: client.opts.Transport.APIKey,
		Message: mailchimpMessage{
			FromName:  client.opts.From.Name,
			FromEmail: client.opts.From.Email,
			Subject:   msg.Subject,
			HTML:      msg.HTML,
			Text:      msg.Text,
			To: []mailchimpRecipient{
				{Name: msg.To.Name, Email: msg.To.Email},
			},
		},
	}

	if err = json.NewEncoder(&buf).Encode(&reqBody); err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, &buf)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer multierr.AppendInvoke(&err, multierr.Close(res.Body))
	resBody := new(mailchimpResponse)
	if err = json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return
	}

	if resBody.Status != "sent" {
		err = fmt.Errorf("received invalid status %q: %q", resBody.Status, resBody.RejectReason)
		return
	}

	return
}

func ConfigureMailchimp(provider config.Provider) (*ClientOptions[MailchimpTransport], error) {
	opts := new(ClientOptions[MailchimpTransport])
	if err := provider.Get("email.mailchimp").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure mailchimp options: %w", err)
	}

	return opts, nil
}

func NewMailchimpClient(opts *ClientOptions[MailchimpTransport]) *MailchimpClient {
	return &MailchimpClient{opts: opts}
}
