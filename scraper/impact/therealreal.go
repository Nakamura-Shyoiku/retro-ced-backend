package impact

import (
	"compress/gzip"
	"encoding/csv"
	"io"

	"net/http"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
)

const (
	impactFTPAddress                    = "products.impactradius.com:21"
	impactFTPUser                       = "ps-ftp_2071858"
	impactFTPPassword                   = "5#xk}hQd2b"
	impactFTPTheRealRealCatalogFilePath = "/The-RealReal/TheRealReal-Product-Catalog_CUSTOM.csv.gz"
)

// CrawlTheRealReal will start 'The Real Real' crawler / FTP processor.
func CrawlTheRealReal(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
	go crawlTheRealReal()
}

/* There is a number of top-level categories that we are ignoring (completely).
    - Art
	- Beauty
	- Home
	- Kids
	- Men
*/

func crawlTheRealReal() {

	if !shouldCrawlTheRealReal() {
		log.Info("theRealReal crawl is disabled - stopping")
		return
	}

	log.Info("starting Impact crawl")

	conn, err := ftp.Dial(impactFTPAddress)
	if err != nil {
		log.WithError(err).Error("could not connect to impact FTP")
		return
	}

	defer conn.Quit()

	err = conn.Login(impactFTPUser, impactFTPPassword)
	if err != nil {
		log.WithError(err).Error("could not login to impact FTP")
		return
	}

	catalog, err := conn.Retr(impactFTPTheRealRealCatalogFilePath)
	if err != nil {
		log.WithError(err).Error("could not open impact catalog file")
		return
	}

	defer catalog.Close()

	log.Info("opened FTP catalog file")

	gz, err := gzip.NewReader(catalog)
	if err != nil {
		log.WithError(err).Error("could not open gzip reader")
		return
	}

	cr := csv.NewReader(gz)

	header, err := cr.Read()
	if err != nil {
		log.WithError(err).Error("could not read file header")
		return

	}

	err = validateHeader(header)
	if err != nil {
		log.WithError(err).Error("file format changed")
		return
	}

	log.Info("validated Impact CSV file format")

	totalSaved := 0

	for i := 0; ; i++ {
		fields, err := cr.Read()
		if err != nil {

			if err == io.EOF {
				break
			}

			log.WithError(err).Error("could not read Impact CSV file")
		}

		// didn't see such data so far
		if fields[deletedFieldNo] != "0" {
			continue
		}

		category, subcategory, keep := extractCategory(fields[categoryFieldNo])
		if !keep {
			continue
		}

		productURL := fields[urlFieldNo]

		prod := models.ProductRecord{
			Guid:        uuid.New().String(),
			SiteId:      12, // TODO: site IDs to const
			UrlId:       0,  // TODO: create URL Guid,
			Category:    category,
			SubCategory: subcategory,
			Brand:       strings.ToLower(fields[brandFieldNo]),
			Title:       fields[nameFieldNo],
			Description: fields[descriptionFieldNo],
			Img:         fields[imageFieldNo],
			Color:       fields[colorFieldNo],
			Size:        fields[sizeFieldNo],
			ItemNumber:  fields[skuFieldNo],
			ProductURL:  &productURL,
			Approved:    true,
		}

		if prod.Category == "shoes" {
			prod.ShoeSize = prod.Size
		}

		// TODO: price as integer - how is it used elsewhere?
		if fields[currencyFieldNo] == "USD" {
			price, _ := strconv.ParseFloat(fields[priceFieldNo], 64)
			prod.RetailPrice = int64(price)

			salePrice, _ := strconv.ParseFloat(fields[salePriceFieldNo], 64)
			prod.Price = int64(salePrice)
		}

		// effectively - if we don't have a saleprice in the original file
		if prod.Price == 0 {
			prod.Price = prod.RetailPrice
		}

		err = models.GetDBv2().Create(&prod).Error
		if err != nil {
			log.WithError(err).Error("could not save impact products to db")
			return
		}

		totalSaved++
	}

	log.WithField("total_saved", totalSaved).Info("completed Impact FTP crawl")
}

func validateHeader(header []string) error {

	expected := []string{
		"sku",
		"name",
		"brand",
		"description",
		"productid",
		"product_id",
		"programname",
		"currency",
		"price",
		"producturl",
		"imageurl",
		"category",
		"saleprice",
		"color",
		"size",
		"instock",
		"alternateimageurl",
		"gender",
		"deleted",
	}

	if len(header) != len(expected) {
		return errors.Errorf("unexpected number of fields (file columns: %+v)", header)
	}

	for i := 0; i < len(expected); i++ {
		if header[i] != expected[i] {
			return errors.Errorf("expected column: '%v' got: '%v'", expected[i], header[i])
		}
	}

	return nil
}

// extractCategory returns product category, subcategory, and boolean - indicating should the
// record be processed or not
func extractCategory(text string) (string, string, bool) {

	fields := strings.Split(text, "|")
	if len(fields) == 0 {
		return "", "", false
	}

	// handle this special case first
	if len(fields) == 3 && fields[0] == "Watches" && strings.TrimSpace(fields[1]) == "Women's Fine Watches" {
		return "accessories", "watches", true
	}

	// immediately discard all uninteresting categories
	if fields[0] != "Women" && fields[0] != "Jewelry" {
		return "", "", false
	}

	// e.g.: Women| Clothing| Dresses| Sleeveless
	if len(fields) >= 3 {

		// e.g. Jewelry| Necklaces| Pendant Necklace| Pendant
		//		return category == accessories, subcategory == jewelry
		if fields[0] == "Jewelry" {
			return "accessories", strings.ToLower(fields[0]), true
		}

		// fields[0] == "Women"

		category := strings.TrimSpace(fields[1])
		if category == "" {
			return "", "", false
		}

		// tweak category for handbags and tops
		if category == "Handbags" {
			category = "bags"
		} else if category == "Tops" {
			category = "clothing"
		}

		return strings.ToLower(category), strings.TrimSpace(strings.ToLower(fields[2])), true
	}

	// not really commonplace
	// discarded stuff at the moment: 'Women||' and 'Jewelry||' - we don't know how
	// to categorize these
	return "", "", false
}

func shouldCrawlTheRealReal() bool {

	var site models.Site
	err := models.GetDBv2().Where("id = 12").First(&site).Error
	if err != nil {
		log.WithError(err).Warn("could not determine if 'TheRealReal' site is disabled")
		return true
	}

	return site.Active
}
