package amorevintage

import (
	"time"

	"github.com/ulventech/retro-ced-backend/utils/client"
)

// Crawler will provide the amorevintage crawling functionality
type Crawler struct {
	cli     *client.Client
	sleep   time.Duration
	maxPage int
	siteID  int

	jpyToUSDConversionRate float64
}

type product struct {
	Name          string            `json:"name,omitempty"`
	Description   string            `json:"description,omitempty"`
	SKU           string            `json:"sku,omitempty"`
	URL           string            `json:"url,omitempty"`
	Category      string            `json:"category,omitempty"`
	ImageURL      string            `json:"image_url,omitempty"`
	Price         string            `json:"price,omitempty"`
	PriceCurrency string            `json:"price_currency,omitempty"`
	Features      map[string]string `json:"features,omitempty"`
}

type searchPage struct {
	Products          []productHTMLInfo         `goquery:"div.grid__item.grid-product"` // little useful info here, but not bad to verify products are lined up correctly
	ProductJSONDocs   []string                  `goquery:"div.product-section:has(script) > script,html"`
	DescriptionTables []productDescriptionTable `goquery:"div#item_details > table.item_info"`
}

type productHTMLInfo struct {
	ProductURL string `goquery:"div.grid-product__content > a,[href]"`
	Title      string `goquery:"div.grid-product__content img.grid-product__image.lazyload,[alt]"`
}

type productDescriptionTable struct {
	ProductName string   `goquery:"tbody > tr:first-child > th.name,text"`
	Name        []string `goquery:"tbody > tr > th:not(th.name),text"`
	Value       []string `goquery:"tbody > tr > td,text"`
}

type productSnippet struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	URL         string `json:"url"`
	Image       struct {
		URL string `json:"url"`
	} `json:"image"`
	Offer struct {
		Price         string `json:"price"`
		PriceCurrency string `json:"priceCurrency"`
	} `json:"offers"`
}
