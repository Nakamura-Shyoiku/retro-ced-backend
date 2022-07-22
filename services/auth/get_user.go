package auth

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type GetUser struct {
	guid		string
	user    models.User
}

func NewGetUser(guid string) *GetUser {
	n := new(GetUser)
	n.guid = guid
	return n
}

func (n *GetUser) Data() models.User {
	return n.user
}

func (n *GetUser) Do() (err error) {
	if err = n.getUser(); err != nil {
		return err
	}

	return nil
}

func (n *GetUser) getUser() (err error) {
	n.user, err = models.GetUserByGUID(n.guid)
	if err != nil {
		return err
	}

	return nil
}