// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package sendgrid

import (
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/email"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/email/templates"
	"github.com/TheThingsNetwork/ttn/pkg/log"
	"github.com/jaytaylor/html2text"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var _ email.Provider = new(SendGrid)

const (
	defaultFromName  = ""
	defaultFromEmail = "noreply@identityserver.ttn"
)

// SendGrid is the type that implements SendGrid as email provider.
type SendGrid struct {
	logger      log.Interface
	client      *sendgrid.Client
	fromEmail   *mail.Email
	renderer    templates.Renderer
	sandboxMode bool
}

// SendGridOpt is the type of functions that configure the provider.
type SendGridOpt func(*SendGrid)

// SandoxMode sets the sandbox mode for testing purposes.
func SandboxMode(enabled bool) SendGridOpt {
	return func(s *SendGrid) {
		s.sandboxMode = enabled
	}
}

// SenderAddress sets the given address as from email address.
func SenderAddress(name, address string) SendGridOpt {
	return func(s *SendGrid) {
		s.fromEmail = mail.NewEmail(name, address)
	}
}

// TemplateRenderer sets the given template renderer and therefore replaces
// the default templates.DefaultRenderer.
func TemplateRenderer(renderer templates.Renderer) SendGridOpt {
	return func(s *SendGrid) {
		s.renderer = renderer
	}
}

// New creates a SendGrid email provider.
func New(logger log.Interface, apiKey string, opts ...SendGridOpt) *SendGrid {
	provider := &SendGrid{
		logger:    logger.WithField("provider", "SendGrid"),
		client:    sendgrid.NewSendClient(apiKey),
		fromEmail: mail.NewEmail(defaultFromName, defaultFromEmail),
		renderer:  templates.DefaultRenderer{},
	}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

// Send sends an email to recipient using the provided template along with the data.
func (s *SendGrid) Send(recipient string, tmpl *templates.Template, data interface{}) error {
	message, err := s.buildEmail(recipient, tmpl, data)
	if err != nil {
		return err
	}

	logger := s.logger.WithFields(log.Fields(
		"recipient", recipient,
		"template.name", tmpl.Name,
	))

	logger.Info("Sending email ...")

	response, err := s.client.Send(message)

	if err != nil {
		logger.WithFields(log.Fields(
			"response.status_code", response.StatusCode,
			"response.body", response.Body,
			"template.data", data,
		)).WithError(err).Error("Failed to send email")

		return err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		logger.WithFields(log.Fields(
			"response.status_code", response.StatusCode,
			"response.body", response.Body,
			"template.data", data,
		)).Error("Failed to send email")

		return errors.Errorf("Failed to send email. Status code `%d`", response.StatusCode)
	}

	logger.Info("Email successfully sent")

	return nil
}

// buildEmail builds the email that will be sent using the underlying SendGrid client.
func (s *SendGrid) buildEmail(to string, tmpl *templates.Template, data interface{}) (*mail.SGMailV3, error) {
	subject, content, err := s.renderer.Render(tmpl, data)
	if err != nil {
		return nil, err
	}

	text, err := html2text.FromString(content, html2text.Options{PrettyTables: true})
	if err != nil {
		return nil, err
	}

	message := mail.NewV3MailInit(
		s.fromEmail,
		subject,
		mail.NewEmail("", to),
		mail.NewContent("text/html", content),
		mail.NewContent("text/plain", text),
	)

	if s.sandboxMode {
		settings := mail.NewMailSettings()
		settings.SetSandboxMode(mail.NewSetting(true))

		message = message.SetMailSettings(settings)
	}

	return message, nil
}
