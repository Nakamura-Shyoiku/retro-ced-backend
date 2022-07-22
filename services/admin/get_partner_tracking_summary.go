package admin

import (
	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

// GetPartnerTrackingSummary service
type GetPartnerTrackingSummary struct {
	userGUID        string
	user            models.User
	site            *models.Site
	trackingSummary models.Tracking
}

// NewGetPartnerTrackingSummary instance
func NewGetPartnerTrackingSummary(userGUID string) *GetPartnerTrackingSummary {
	n := new(GetPartnerTrackingSummary)
	n.userGUID = userGUID
	return n
}

// Data returns all users
func (g *GetPartnerTrackingSummary) Data() models.Tracking {
	return g.trackingSummary
}

// Do task
func (g *GetPartnerTrackingSummary) Do() (err error) {
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

func (g *GetPartnerTrackingSummary) getUser() (err error) {
	g.user, err = models.GetUserByGUID(g.userGUID)
	if err != nil {
		log.WithError(err).Error("failed to get site name")
		return err
	}

	return nil
}

func (g *GetPartnerTrackingSummary) getSite() (err error) {
	g.site, err = new(models.Site).GetByID(g.user.PartnerSiteID)
	if err != nil {
		log.WithError(err).Warn("failed to get site by id")
		return err
	}

	return nil
}

func (g *GetPartnerTrackingSummary) getTrackingSummary() (err error) {
	g.trackingSummary, err = new(models.Tracking).GetTrackingSummary(g.site.Name)
	if err != nil {
		log.WithError(err).Error("failed to get site name")
		return err
	}

	return nil
}
