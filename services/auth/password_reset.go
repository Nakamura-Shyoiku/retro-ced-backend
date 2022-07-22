package auth

import (
	"errors"

	"github.com/apex/log"
	"github.com/badoux/checkmail"
	"github.com/spf13/viper"
	"github.com/thedanielforum/randtool"

	"github.com/ulventech/retro-ced-backend/email"
	"github.com/ulventech/retro-ced-backend/models"
)

// PasswordReset service
type PasswordReset struct {
	email string
	token string
}

// NewPasswordReset instance
func NewPasswordReset(email string) *PasswordReset {

	return &PasswordReset{
		email: email,
	}
}

// Do tasks
func (p *PasswordReset) Do() (err error) {
	if err = p.validate(); err != nil {
		return err
	}

	if err = p.generateToken(); err != nil {
		return err
	}

	if err = p.sendEmail(); err != nil {
		return err
	}

	return nil
}

func (p *PasswordReset) validate() (err error) {
	if p.email == "" {
		err = errors.New("email is required")
		log.Warn(err.Error())
		return err
	}

	if err = checkmail.ValidateFormat(p.email); err != nil {
		log.WithError(err).Warn("email format error")
		return err
	}

	return nil
}

func (p *PasswordReset) generateToken() (err error) {
	p.token, err = randtool.GenStr(60)
	if err != nil {
		log.WithError(err).Error("failed to generate rand string")
		return err
	}

	return new(models.User).UpdatePasswordReset(
		p.token,
		p.email,
	)
}

func (p *PasswordReset) sendEmail() (err error) {

	hostname := "www.retroced.com"
	if viper.GetBool("app.beta") {
		hostname = "beta.retroced.com"
	}

	data := email.Data{
		"Token":    p.token,
		"Email":    p.email,
		"Hostname": hostname,
	}

	err = email.NewEmail().SendPasswordResetToken(p.email, data)

	if err != nil {
		// Just log and continue execution
		log.WithError(err).Error("failed to send SendPasswordResetToken")
	}
	return nil
}
