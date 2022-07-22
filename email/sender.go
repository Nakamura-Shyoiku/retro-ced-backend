package email

import (
	"bytes"
	"html/template"

	"github.com/apex/log"
	"github.com/spf13/viper"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// Email client
type Email struct {
	client mailgun.Mailgun
}

// Data that is passes to go template parser {{ .Name }}
type Data map[string]interface{}

// NewEmail instance
func NewEmail() *Email {

	return &Email{
		client: mailgun.NewMailgun(
			viper.GetString("mailgun.domain"),
			viper.GetString("mailgun.secret"),
			viper.GetString("mailgun.key"),
		),
	}
}

// SendPasswordResetToken sends password reset token
func (e *Email) SendPasswordResetToken(recipient string, data Data) error {
	// process email template
	t := template.Must(template.New("password_reset.html").ParseFiles("./email/templates/password_reset.html"))
	body := new(bytes.Buffer)
	err := t.Execute(body, data)
	if err != nil {
		log.WithError(err).Error("failed to parse template")
		return err
	}

	message := e.client.NewMessage(
		viper.GetString("mailgun.sender"),
		viper.GetString("email_subjects.password_reset"),
		"Please view this email with a HTML compatible email client",
		recipient,
	)

	message.SetHtml(body.String())

	resp, id, err := e.client.Send(message)
	if err != nil {
		log.WithError(err).Error("failed to send email")
		return err
	}

	log.WithField("Guid", id).
		WithField("resp", resp).
		Info("sent password reset email")

	return nil
}

// SendContactUsEmail sends email to site admin
func (e *Email) SendContactUsEmail(sender string, data Data) error {
	// process email template
	t := template.Must(template.New("contact_us.html").ParseFiles("./email/templates/contact_us.html"))
	body := new(bytes.Buffer)
	err := t.Execute(body, data)
	if err != nil {
		log.WithError(err).Error("failed to parse template")
		return err
	}

	message := e.client.NewMessage(
		sender,
		viper.GetString("email_subjects.contact_us"),
		"Please view this email with a HTML compatible email client",
		viper.GetString("mailgun.admin"),
	)
	message.SetHtml(body.String())
	resp, id, err := e.client.Send(message)
	if err != nil {
		log.WithError(err).Error("failed to send email")
		return err
	}

	log.Infof("Guid: %s Resp: %s", id, resp)
	return nil
}

// SendPartnerApplication sends application email to site admin
func (e *Email) SendPartnerApplication(sender string, data Data) error {
	// process email template
	t := template.Must(template.New("partner_application.html").ParseFiles("./email/templates/partner_application.html"))
	body := new(bytes.Buffer)
	err := t.Execute(body, data)
	if err != nil {
		log.WithError(err).Error("failed to parse template")
		return err
	}

	message := e.client.NewMessage(
		sender,
		viper.GetString("email_subjects.partner_application"),
		"Please view this email with a HTML compatible email client",
		viper.GetString("mailgun.admin"),
	)
	message.SetHtml(body.String())
	resp, id, err := e.client.Send(message)
	if err != nil {
		log.WithError(err).Error("failed to send email")
		return err
	}

	log.Infof("Guid: %s Resp: %s", id, resp)
	return nil
}
