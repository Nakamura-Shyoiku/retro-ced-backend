package admin

import (
	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
)

// DeleteSite service
type DeleteSite struct {
	SiteID   int64
	products []models.Product
}

// NewDeleteSite instance
func NewDeleteSite(siteID int64) *DeleteSite {
	return &DeleteSite{
		SiteID: siteID,
	}
}

// Do tasks
func (ds *DeleteSite) Do() (err error) {
	if err = ds.validate(); err != nil {
		return err
	}

	err = ds.deleteProducts()
	if err != nil {
		log.WithError(err).
			WithField("site_id", ds.SiteID).
			Info("failed to delete products related to the site")
	}

	if err = ds.deleteSite(); err != nil {
		return err
	}

	return nil
}

func (ds *DeleteSite) validate() (err error) {
	if ds.SiteID <= 0 {
		err = errors.New("'site_id' is required")
		log.Info(err.Error())
		return err
	}

	return nil
}

func (ds *DeleteSite) deleteSite() error {
	return new(models.Site).Delete(ds.SiteID)
}

func (ds *DeleteSite) deleteProducts() error {

	stmt, err := models.GetDB().Prepare(`
		ALTER TABLE Products DELETE WHERE site_id = ?
	`)
	if err != nil {
		log.WithError(err).Error("failed prepare delete Site statement")
		return err
	}
	_, err = stmt.Exec(ds.SiteID)
	if err != nil {
		log.WithError(err).Error("failed to run exec on delete Site")
		return err
	}

	// deleter := sqlbuilder.NewDeleteBuilder()
	// deleter.DeleteFrom("Products")
	// deleter.Where(deleter.Equal("site_id", ds.SiteID))

	// query, params := deleter.Build()

	// _, err := models.GetDB().Query(query, params...)
	// if err != nil {
	// 	return errors.Wrap(err, "could not remove products")
	// }

	return nil
}
