package services

import (
	"errors"
	
	"github.com/apex/log"
	
	"github.com/ulventech/retro-ced-backend/email"
)

// ContactUs service
type ContactUs struct {
	senderEmail string
	name        string
	phone       string
	message     string
}

// NewContactUs instance
func NewContactUs(senderEmail, name, phone, message string) *ContactUs {
	return &ContactUs{
		senderEmail,
		name,
		phone,
		message,
	}
}

// Do tasks
func (cu *ContactUs) Do() (err error) {
	if err = cu.validate(); err != nil {
		return err
	}

	if err = cu.sendMessage(); err != nil {
		return err
	}

	return nil
}

func (cu *ContactUs) validate() (err error) {
	if cu.senderEmail == "" {
		err = errors.New("sender email is required")
		log.Warn(err.Error())
		return err
	}

	if cu.message == "" {
		err = errors.New("message is required")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (cu *ContactUs) sendMessage() error {
	err := email.NewEmail().SendContactUsEmail(cu.senderEmail, email.Data{
		"Name":    cu.name,
		"Email":   cu.senderEmail,
		"Phone":   cu.phone,
		"Message": cu.message,
	})
	if err != nil {
		log.WithError(err).Error("failed to send contact us form")
		return err
	}

	return nil
}
