package admin

import (
	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

type GetSites struct {
	sites []models.Site
}

func NewGetSites() *GetSites {
	return new(GetSites)
}

func (g *GetSites) Data() []models.Site {
	return g.sites
}

func (g *GetSites) Do() error {
	return g.getSites()
}

func (g *GetSites) getSites() (err error) {

	g.sites, err = models.GetAllSites()
	if err != nil {
		log.WithError(err).Error("failed to get sites")
		return err
	}

	return nil
}
