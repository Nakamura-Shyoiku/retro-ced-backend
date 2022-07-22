package auth

import (
	"errors"
	"fmt"
	"github.com/ulventech/retro-ced-backend/utils/env"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"github.com/apex/log"
	"github.com/badoux/checkmail"
	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils"
)

// Register service
type Register struct {
	email string
	pass  string
	user  models.User
}

// NewRegister instance
func NewRegister(email, pass string) *Register {
	n := new(Register)
	n.email = email
	n.pass = pass
	return n
}

// Do tasks
func (r *Register) Do() (err error) {
	if err = r.validate(); err != nil {
		return err
	}

	if err = r.checkIfUserExists(); err != nil {
		return err
	}

	if err = r.register(); err != nil {
		return err
	}

	if err = r.addToMailchimp(); err != nil {
		return err
	}

	return nil
}

// Data returns jwt token
func (r *Register) Data() string {
	return utils.NewToken(
		fmt.Sprintf("%d", r.user.Id),
		r.user.Email,
		r.user.Username,
		r.user.FirstName,
		r.user.LastName,
		r.user.FbID,
		r.user.ACL,
	)
}

func (r *Register) validate() (err error) {
	if r.email == "" {
		err = errors.New("email is required")
		log.Warn(err.Error())
		return err
	}

	if err = checkmail.ValidateFormat(r.email); err != nil {
		log.WithError(err).Warn("email format error")
		return err
	}

	if r.pass == "" {
		err = errors.New("password is required")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (r *Register) checkIfUserExists() (err error) {
	r.user, err = models.GetUserByEmail(r.email)
	if err != nil {
		log.WithError(err).Error("error while fetching user")
		return err
	}

	if r.user.Id > 0 {
		err = errors.New("there is already a user with this email")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (r *Register) register() error {
	// hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(r.pass),
		viper.GetInt("app.bcrypt_cost"),
	)
	if err != nil {
		return err
	}

	// populate user struct
	user := &models.User{
		Email:    r.email,
		Password: fmt.Sprintf("%s", hash),
	}

	_, err = user.Create()
	if err != nil {
		return err
	}

	r.user, err = models.GetUserByEmail(r.email)
	if err != nil {
		log.WithError(err).Error("error while fetching user")
		return err
	}

	return nil
}

func (r *Register) addToMailchimp() (err error) {
	if !env.IsProd() {
		return nil
	}

	err = mailchimp.SetKey(viper.GetString("mailchimp.key"))
	if err != nil {
		return err
	}

	params := &members.NewParams{
		EmailAddress: r.email,
		Status:       members.StatusSubscribed,
	}

	_, err = members.New(viper.GetString("mailchimp.signups"), params)
	if err != nil {
		log.WithError(err).Error("failed to add email to list")
		return err
	}

	return nil
}
