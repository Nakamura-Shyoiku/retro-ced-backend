package admin

import (
	"io/ioutil"

	"github.com/apex/log"
)

// ClickTrackingSummary service
type ClickTrackingSummary struct {
	page string
}

// ClickTrackingSummary instance
func NewClickTrackingSummary() *ClickTrackingSummary {
	return new(ClickTrackingSummary)
}

func (c *ClickTrackingSummary) Data() string {
	return c.page
}

// Do tasks
func (c *ClickTrackingSummary) Do() (err error) {
	return c.readFile()
}

func (c *ClickTrackingSummary) readFile() error {
	htmlData, err := ioutil.ReadFile("./static/admin/click_tracking_summary.html")
	if err != nil {
		log.WithError(err).Error("failed to read click_tracking_summary.html")
		return err
	}
	c.page = string(htmlData)
	return nil
}
