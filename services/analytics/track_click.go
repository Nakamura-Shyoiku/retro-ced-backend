package services

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type TrackClick struct {
	link string
}

func NewTrackClick(link string) *TrackClick {
	n := new(TrackClick)
	n.link = link
	return n
}

func (t *TrackClick) Do() (err error) {
	return t.saveClick()
}

func (t *TrackClick) saveClick() error {
	return new(models.Click).Create(t.link)
}
