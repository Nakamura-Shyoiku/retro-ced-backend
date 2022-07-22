package auth

import (
	"errors"
	"fmt"
	
	"github.com/badoux/checkmail"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

// SetPassword service
type SetPassword struct {
	token string
	email string
	pass  string
	user  models.User
}

// NewSetPassword instance
func NewSetPassword(email, token, pass string) *SetPassword {
	n := new(SetPassword)
	n.token = token
	n.email = email
	n.pass = pass
	return n
}

// Do tasks
func (s *SetPassword) Do() (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if err = s.getUser(); err != nil {
		return err
	}

	if err = s.updateUser(); err != nil {
		return err
	}

	return nil
}

func (s *SetPassword) validate() (err error) {
	if s.token == "" {
		err = errors.New("token is required")
		log.Warn(err.Error())
		return err
	}

	if s.pass == "" {
		err = errors.New("password is required")
		log.Warn(err.Error())
		return err
	}

	if s.email == "" {
		err = errors.New("email is required")
		log.Warn(err.Error())
		return err
	}

	if err = checkmail.ValidateFormat(s.email); err != nil {
		log.WithError(err).Warn("email format error")
		return err
	}

	return nil
}

func (s *SetPassword) getUser() (err error) {
	s.user, err = models.GetUserByEmail(s.email)
	if err != nil {
		log.WithError(err).Error("failed to get user")
		return err
	}

	// check for valid user and token
	if s.user.Id <= 0 {
		err = errors.New("invalid email or token, please try requesting a new email")
		log.Warn(err.Error())
		return err
	}

	// check that current user token is not empty
	if s.user.PasswordReset == "" {
		err = errors.New("invalid email or token, please try requesting a new email")
		log.Warn(err.Error())
		return err
	}

	if s.user.PasswordReset != s.token {
		err = errors.New("invalid email or token, please try requesting a new email")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (s *SetPassword) updateUser() error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(s.pass),
		viper.GetInt("app.bcrypt_cost"),
	)
	if err != nil {
		log.WithError(err).Error("failed to encrypt password")
		return err
	}

	// remove token
	user := new(models.User)
	if err = user.UpdatePasswordReset("", s.user.Email); err != nil {
		log.WithError(err).Error("failed to update token")
		return err
	}

	// set new pass
	err = user.UpdateUserPass(
		fmt.Sprintf("%s", hash),
		s.user.Email,
	)
	if err != nil {
		log.WithError(err).Error("failed to update password")
		return err
	}

	return nil
}
