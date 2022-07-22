package admin

import (
	"errors"

	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

type UpdateApprovedStatus struct {
	Guid     []string `json:"guid"`
	Approved bool     `json:"approved"`
	products []models.Product
}

func NewUpdateApprovedStatus(data UpdateApprovedStatus) *UpdateApprovedStatus {
	n := new(UpdateApprovedStatus)
	n.Guid = data.Guid
	n.Approved = data.Approved
	return n
}

func (u *UpdateApprovedStatus) Do() (err error) {
	if err = u.validate(); err != nil {
		return err
	}

	if err = u.getProduct(); err != nil {
		return err
	}

	if err = u.updateApprovedStatus(); err != nil {
		return err
	}

	return nil
}

func (u *UpdateApprovedStatus) validate() (err error) {
	if len(u.Guid) <= 0 {
		return errors.New("id is required")
	}
	return nil
}

// func (u *UpdateApprovedStatus) getProduct() (err error) {
// 	product, err = models.GetProductByGuid(int64(u.Guid))
// 	if err != nil {
// 		log.WithError(err).Error("failed to getProduct")
// 		return err
// 	}

// 	return nil
// }

func (u *UpdateApprovedStatus) getProduct() (err error) {
	u.products, err = models.GetProductsByGuid(u.Guid)
	if err != nil {
		log.WithError(err).Error("failed to getProduct")
		return err
	}
	return nil

}

func (u *UpdateApprovedStatus) updateApprovedStatus() (err error) {
	log.Infof("Will update: %s", u.Guid)
	err = new(models.Product).UpdateApprovedStatus(u.Guid, u.Approved)
	if err != nil {
		log.WithError(err).Error("failed to UpdateApprovedStatus")
		return err
	}

	return nil
}
