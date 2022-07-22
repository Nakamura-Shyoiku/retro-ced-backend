package amorevintage

import (
	"encoding/json"
	"html"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// this error indicates that the arrays of product info - HTML, JSON and table data -
// we extracted from the search page are of unequal lengths. This means something went
// wrong and we cannot pair the sections to get the complete product data.
var errSectionLengthMismatch = errors.Errorf("product sections length mismatch")

// this error indicates that a specific product table extracted - th/td pairs -
// is of unequal length. This means we cannot pair keys and values for the specific product.
var errTableLengthMismatch = errors.Errorf("table key-value length mismatch")

// transform will validate and translate raw HTML search page to a list of products
func (p *searchPage) transform() ([]product, error) {

	// all three arrays must have the same length
	if len(p.Products) != len(p.ProductJSONDocs) {

		log.WithField("products_html_count", len(p.Products)).
			WithField("products_json_count", len(p.ProductJSONDocs)).
			WithField("description_tables_count", len(p.DescriptionTables)).
			Error("product sections length mismatch")

		return nil, errSectionLengthMismatch
	}

	var descriptionTablesMismatch bool

	if len(p.DescriptionTables) != len(p.Products) {
		log.WithField("products", len(p.Products)).WithField("description_tables", len(p.DescriptionTables)).Warn("more description tables than required - possible duplicate")
		descriptionTablesMismatch = true
	}

	descriptionMap := make(map[string]productDescriptionTable)
	for _, desc := range p.DescriptionTables {
		name := strings.TrimSpace(desc.ProductName)

		log.WithField("name", name).Debug("description map table set")
		descriptionMap[name] = desc
	}

	products := make([]product, len(p.Products))

	for i := 0; i < len(p.Products); i++ {

		// unmarshal JSON snippet
		var jsonDoc productSnippet
		err := json.Unmarshal([]byte(html.UnescapeString(p.ProductJSONDocs[i])), &jsonDoc)
		if err != nil {

			log.WithField("index", i).WithField("json_string", p.ProductJSONDocs[i]).Error("could not parse product JSON snippet")

			return nil, errors.Wrap(err, "could not unmarshal product JSON snippet")
		}

		if jsonDoc.Name != p.Products[i].Title {
			log.WithField("json_name", jsonDoc.Name).
				WithField("html_title", p.Products[i].Title).
				Warn("product names not matching")
		}

		products[i] = product{
			Name:          jsonDoc.Name,
			Description:   jsonDoc.Description,
			SKU:           jsonDoc.SKU,
			URL:           jsonDoc.URL,
			Category:      jsonDoc.Category,
			ImageURL:      jsonDoc.Image.URL,
			Price:         jsonDoc.Offer.Price,
			PriceCurrency: jsonDoc.Offer.PriceCurrency,
		}

		var productParameters productDescriptionTable

		if !descriptionTablesMismatch {
			productParameters = p.DescriptionTables[i]
		} else {

			// process table data - key/value pairs
			var found bool
			productParameters, found = descriptionMap[strings.TrimSpace(jsonDoc.Name)]
			if !found {

				log.WithField("product_name", p.Products[i].Title).
					WithField("json_name", jsonDoc.Name).
					Error("could not locate description table for product")

				continue
			}
		}

		err = productParameters.valid()
		if err != nil {

			log.WithField("index", i).
				WithField("names", productParameters.Name).
				WithField("values", productParameters.Value).
				Error("invalid description table")

			return nil, err
		}

		products[i].Features = productParameters.GetMapFromTable()
	}

	return products, nil
}

func (t productDescriptionTable) valid() error {

	if len(t.Name) != len(t.Value) {
		return errTableLengthMismatch
	}

	return nil
}

func (t productDescriptionTable) GetMapFromTable() map[string]string {

	params := make(map[string]string)

	for j := 0; j < len(t.Name); j++ {
		name := strings.TrimSpace(t.Name[j])
		value := strings.TrimSpace(t.Value[j])

		params[name] = value
	}

	return params
}

func (p product) getColor() string {

	field := p.Features["Color / material"]
	if field == "" {
		return ""
	}

	fields := strings.Split(field, "/")
	if len(fields) == 0 {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(fields[0]))
}

func (p product) getCondition() string {
	return strings.ToLower(strings.TrimSpace(p.Features["Outside"]))
}

func (p product) getSize() string {
	return strings.TrimSpace(p.Features["Size(cm)"])
}

func (p product) getAccessories() string {
	return strings.TrimSpace(p.Features["Accessories"])
}

func (c *Crawler) getPriceInUSD(p product) float64 {

	if p.Price == "" || p.PriceCurrency == "" {
		return 0
	}

	price, _ := strconv.ParseFloat(p.Price, 64)
	if price == 0 {
		return 0
	}

	if p.PriceCurrency == "USD" {
		return price
	}

	if p.PriceCurrency == "JPY" {
		return c.jpyToUSDConversionRate * price
	}

	log.WithField("currency", p.PriceCurrency).Warn("found unexpected product currency")

	return 0
}
