package admin

import (
	"github.com/apex/log"
	
	"github.com/ulventech/retro-ced-backend/models"
)

// GetPartnerTrackingSummary service
type GetProductTrackingSummary struct {
	userGUID        string
	fromDate        string
	toDate          string
	user            models.User
	site            *models.Site
	trackingSummary []models.Tracking
}

// NewGetPartnerTrackingSummary instance
func NewGetProductTrackingSummary(userGUID string, fromDate string, toDate string) *GetProductTrackingSummary {
	n := new(GetProductTrackingSummary)
	n.userGUID = userGUID
	n.fromDate = fromDate
	n.toDate = toDate
	return n
}

// Data returns all users
func (g *GetProductTrackingSummary) Data() []models.Tracking {
	return g.trackingSummary
}

// Do task
func (g *GetProductTrackingSummary) Do() (err error) {
	if err = g.getUser(); err != nil {
		return err
	}

	if err = g.getSite(); err != nil {
		return err
	}

	if err = g.getTrackingSummary(); err != nil {
		return err
	}

	return nil
}

func (g *GetProductTrackingSummary) getUser() (err error) {
	g.user, err = models.GetUserByGUID(g.userGUID)
	if err != nil {
		log.WithError(err).Error("failed to get site name")
		return err
	}

	return nil
}

func (g *GetProductTrackingSummary) getSite() (err error) {
	g.site, err = new(models.Site).GetByID(g.user.PartnerSiteID)
	if err != nil {
		log.WithError(err).Warn("failed to get site by id")
		return err
	}

	return nil
}

func (g *GetProductTrackingSummary) getTrackingSummary() (err error) {
	g.trackingSummary, err = new(models.Tracking).GetProductTrackingSummary(g.site.Name, g.fromDate, g.toDate)
	if err != nil {
		log.WithError(err).Error("failed to get summary")
		return err
	}

	return nil
}
