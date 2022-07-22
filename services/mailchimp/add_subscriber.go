package mailchimp

import (
	"github.com/apex/log"
	checkmail "github.com/badoux/checkmail"
	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/spf13/viper"
)

type AddSubscriber struct {
	email string
}

func NewAddSubscriber(email string) *AddSubscriber {
	n := new(AddSubscriber)
	n.email = email
	return n
}

func (a *AddSubscriber) Do() (err error) {
	if err = a.validate(); err != nil {
		return err
	}

	if err = a.add(); err != nil {
		return err
	}

	return nil
}

func (a *AddSubscriber) validate() (err error) {
	if err = checkmail.ValidateFormat(a.email); err != nil {
		log.WithError(err).Info("validate email")
		return err
	}

	return nil
}

func (a *AddSubscriber) add() (err error) {
	err = mailchimp.SetKey(viper.GetString("mailchimp.key"))
	if err != nil {
		return err
	}

	params := &members.NewParams{
		EmailAddress: a.email,
		Status:       members.StatusSubscribed,
	}

	_, err = members.New(viper.GetString("mailchimp.newsletter"), params)
	if err != nil {
		log.WithError(err).Error("failed to add email to list")
		return err
	}

	return nil
}
