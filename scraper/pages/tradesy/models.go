package tradesy

type resultPage struct {
	Items []string `goquery:"div#category-item-grid div.item-tile__image-layout a.item-tile__image-link,[href]"`
}

type adPage struct {
	MetadataJSON  string `goquery:"script#liftigniter-metadata,text"`
	FeaturedImage string `goquery:"div#idp-content div[data-testid=Gallery] div[data-testid=ImageSlider] img,[src]"`

	// NOTE: not currently used
	Images []string `goquery:"div#idp-react-app div[data-testid=Gallery] div[data-testid=ImageSlider] img.lazyload,[data-src]"`
}

// JSON payload describing the ad, delivered within HTML
type adMetadata struct {
	Brand        string  `json:"brand"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Condition    string  `json:"condition"`
	Color        string  `json:"color"`
	Size         string  `json:"size"` // use this one for DB size field
	Measurements string  `json:"measurements"`

	// TODO: havent found any ads with retail price.
	// There's more data in the HTML itself.
}
