package malleries

import (
	"strings"
	"time"

	"github.com/ulventech/retro-ced-backend/utils/client"
)

const (
	malleriesEndpoint = "https://www.malleries.com/items"
)

type resultsPage struct {
	Items []string `goquery:"div#products > div.item > div.thumbnail > a,[href]"`
}

// Crawler will crawl the malleries website
type Crawler struct {
	cli     *client.Client
	maxPage int // TODO: get from config
	sleep   time.Duration
	siteID  uint64
}

type rawAdPage struct {
	Title                  string   `goquery:"h1.text-center"`
	DescriptionSections    []string `goquery:"div[itemprop=\"description\"] > p,text"`
	DescriptionSectionsAlt []string `goquery:"div.description > div[style=\"padding-top:10px;\"]:has(strong),text"`
	ItemModel              string   `goquery:"span.model"`
	Price                  float64  `goquery:"span[itemprop=\"price\"],[content]"`
	PriceCurrency          string   `goquery:"span[itemprop=\"priceCurrency\"],[content]"`
	Images                 []string `goquery:"div#itemImageCarousel div.carousel-inner > div.item.thumbnail.item-image-wrapper > img,[data-full]"`
	BreadCrumbs            []string `goquery:"ol.breadcrumb > li > a,text"`
}

// TODO: perhaps obsolete this and just translate from rawAdPage to models.Product
type adPage struct {
	Title         string
	Description   string
	Brand         string
	ItemModel     string
	Price         float64
	PriceCurrency string
	Images        []string
	Measurements  string
	Color         string
	RetailPrice   string
	Size          string
}

const (
	measurementsPrefix            = "Measurements:"
	measurementsPrefix2           = "Measurement:"
	measurementsInDescriptionText = "See description"
	otherDesignersLabel           = "Other Designers"

	colorPrefix       = "Color:"
	retailPricePrefix = "Retail price:"
	sizePrefix        = "Size:"
)

func (p rawAdPage) adPage() adPage {

	out := adPage{
		Title:         p.Title,
		Description:   strings.Join(p.DescriptionSections, "\n"),
		ItemModel:     p.ItemModel,
		Images:        p.Images,
		Price:         p.Price,
		PriceCurrency: p.PriceCurrency,
	}

	if len(p.BreadCrumbs) >= 3 {
		if !strings.Contains(p.BreadCrumbs[2], otherDesignersLabel) {
			out.Brand = strings.TrimSpace(p.BreadCrumbs[2])
		}
	}

	sections := make([]string, 0, len(p.DescriptionSections))
	sections = append(sections, p.DescriptionSections...)
	sections = append(sections, p.DescriptionSectionsAlt...)

	for i := 0; i < len(sections); i++ {

		section := sections[i]

		// TODO: use regex
		if strings.HasPrefix(section, measurementsPrefix) && !strings.Contains(section, measurementsInDescriptionText) {
			out.Measurements = strings.TrimSpace(strings.TrimPrefix(section, measurementsPrefix))
			continue
		}

		if strings.HasPrefix(section, measurementsPrefix2) && !strings.Contains(section, measurementsInDescriptionText) {
			out.Measurements = strings.TrimSpace(strings.TrimPrefix(section, measurementsPrefix2))
			continue
		}

		if strings.HasPrefix(section, colorPrefix) {
			out.Color = strings.TrimSpace(strings.TrimPrefix(section, colorPrefix))
			continue
		}

		// TODO: get float only
		if strings.HasPrefix(section, retailPricePrefix) {
			out.RetailPrice = strings.TrimSpace(strings.TrimPrefix(section, retailPricePrefix))
			continue
		}

		if strings.HasPrefix(section, sizePrefix) {
			out.Size = strings.TrimSpace(strings.TrimPrefix(section, sizePrefix))
		}
	}

	return out
}
