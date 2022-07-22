package malleries

import (
	"fmt"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

func PaginateBags(s models.Site) {
	log.Infof("paginating: %s", s.Name)
	for p := int64(1); p <= s.MaxPage; p++ {
		// If sleep if configured
		if s.Sleep > 0 {
			log.Warnf("sleeping for %dms", s.Sleep)
			time.Sleep(time.Duration(s.Sleep) * time.Millisecond)
		}
		log.Infof("scraping page %d for product links", p)

		resp, err := soup.Get(fmt.Sprintf("http://www.malleries.com/all.php?page=%d", p))
		if err != nil {
			log.WithError(err).Fatal("failed to fetch")
		}
		d := soup.HTMLParse(resp)
		log.Infof("%+v", d)
	}
}
