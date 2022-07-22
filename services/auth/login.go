package auth

import (
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils"
)

// Login service
type Login struct {
	email string
	pass  string
	user  models.User
}

// NewLogin instance
func NewLogin(email, pass string) *Login {
	n := new(Login)
	n.email = email
	n.pass = pass
	return n
}

// Do tasks
func (l *Login) Do() (err error) {
	if err = l.validate(); err != nil {
		return err
	}

	if err = l.checkIfUserExists(); err != nil {
		return err
	}

	if err = l.checkPassword(); err != nil {
		return err
	}

	return nil
}

// Data returns auth string
func (l *Login) Data() string {
	return utils.NewToken(
		fmt.Sprintf("%d", l.user.Id),
		l.user.Email,
		l.user.Username,
		l.user.FirstName,
		l.user.LastName,
		l.user.FbID,
		l.user.ACL,
	)
}

func (l *Login) validate() (err error) {
	if l.email == "" {
		err = errors.New("email is required")
		log.Warn(err.Error())
		return err
	}

	if err = checkmail.ValidateFormat(l.email); err != nil {
		log.WithError(err).Warn("email format error")
		return err
	}

	if l.pass == "" {
		err = errors.New("password is required")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (l *Login) checkIfUserExists() (err error) {
	l.user, err = models.GetUserByEmail(l.email)
	if err != nil {
		log.WithError(err).Error("error while fetching user")
		return err
	}

	if l.user.Password == "" {
		err = errors.New("please reset your password")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (l *Login) checkPassword() error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(l.user.Password),
		[]byte(l.pass),
	)
	if err != nil {
		err = errors.New("password or email is invalid")
		log.Warn(err.Error())
		return err
	}

	return nil
}

// HasPermission verifies user rights for a specific project
func HasPermission(userID string, permissionID int64) bool {
	exists := true

	err := models.VerifyPermission(userID, permissionID)
	if err != nil {
		exists = false
	}

	return exists
}
