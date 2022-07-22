package admin

import (
	"errors"
	
	"github.com/apex/log"
	
	"github.com/ulventech/retro-ced-backend/models"
)

//GetUsers service
type GetUsers struct {
	searchQuery string
	offset      int
	limit       int
	totalCnt    uint64
	users       []models.User
}

//NewGetUsers instance
func NewGetUsers(offset, limit int, searchQuery string) *GetUsers {
	n := new(GetUsers)
	n.offset = offset
	n.limit = limit
	n.searchQuery = searchQuery
	return n
}

//Data returns all users
func (g *GetUsers) Data() ([]models.User, uint64) {
	return g.users, g.totalCnt
}

//Do task
func (g *GetUsers) Do() (err error) {
	if err = g.validate(); err != nil {
		return err
	}

	if err = g.getUsers(); err != nil {
		return err
	}

	if err = g.getUsersCount(); err != nil {
		return err
	}

	return nil
}

func (g *GetUsers) validate() (err error) {
	if g.offset < 0 {
		err = errors.New("offset can not be negative")
		log.Warn(err.Error())
		return err
	}

	if g.limit < 0 {
		err = errors.New("limit can not be negative")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (g *GetUsers) getUsers() (err error) {
	g.users, err = models.GetUsers(g.searchQuery, g.offset, g.limit)
	if err != nil {
		log.WithError(err).Error("failed to get users")
		return err
	}

	return nil
}

func (g *GetUsers) getUsersCount() (err error) {
	g.totalCnt, err = models.GetUsersCount(g.searchQuery)
	if err != nil {
		log.WithError(err).Error("failed to get users count")
		return err
	}
	return nil
}
