package admin

import (
	"errors"
	
	"github.com/apex/log"
	
	"github.com/ulventech/retro-ced-backend/models"
)

// UpdateUser service
type UpdateUser struct {
	userID int64
	acl    int
	siteID int64
	user   *models.User
}

// NewUpdateUser instance
func NewUpdateUser(userID, acl, siteID int64) *UpdateUser {
	return &UpdateUser{
		userID: userID,
		siteID: siteID,
		acl:    int(acl),
		user:   new(models.User),
	}
}

// Data returns the updated user
func (u *UpdateUser) Data() *models.User {
	return u.user
}

// Do tasks
func (u *UpdateUser) Do() (err error) {
	if err = u.validate(); err != nil {
		return err
	}

	if err = u.update(); err != nil {
		return err
	}

	if err = u.getUpdated(); err != nil {
		return err
	}

	return nil
}

func (u *UpdateUser) validate() (err error) {
	if u.userID <= 0 {
		err = errors.New("userID is invalid")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (u *UpdateUser) update() (err error) {
	err = new(models.User).UpdateUserACL(
		u.userID,
		u.acl,
		u.siteID,
	)
	if err != nil {
		log.WithError(err).Warn("failed to change the users permission level")
		return err
	}

	return nil
}

func (u *UpdateUser) getUpdated() (err error) {
	u.user, err = new(models.User).GetUserByID(u.userID)
	if err != nil {
		log.WithError(err).Warn("failed to get user by id")
		return err
	}

	return nil
}
