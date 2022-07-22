package admin

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type Update struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Sleep    int64  `json:"sleep"`
	MaxPages uint64 `json:"max_page"`
	Active   bool   `json:"active"`
}

func NewUpdate(data Update) *Update {
	n := new(Update)
	n.ID = data.ID
	n.Name = data.Name
	n.Sleep = data.Sleep
	n.MaxPages = data.MaxPages
	n.Active = data.Active
	return n
}

// Do tasks
func (u *Update) Do() (err error) {
	if err = u.validate(); err != nil {
		return err
	}

	if err = u.updateSite(); err != nil {
		return err
	}

	return nil
}

func (u *Update) validate() (err error) {
	return nil
}

func (u *Update) updateSite() error {
	return new(models.Site).Update(
		u.ID,
		u.MaxPages,
		u.Sleep,
		u.Name,
		u.Active,
	)
}
