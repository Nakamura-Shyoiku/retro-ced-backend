package whatgoesaround

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ulventech/retro-ced-backend/models"
)

const (
	endMarker       = "TRL"
	fieldsPerRecord = 38

	// batchSize = 500

	siteID = 7

	fieldID          = 0
	fieldName        = 1
	fieldCategory    = 3
	fieldSubCategory = 4
	fieldURL         = 5
	fieldImage       = 6
	fieldPrice       = 12
	fieldRetailPrice = 13
	fieldBrand       = 20
	fieldCurrency    = 25
	fieldSize        = 30
	fieldColor       = 32
)

var catalogRe = regexp.MustCompile(`42946_\d+_mp.txt.gz`)

func Crawl(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
	go crawl()
}

func crawl() {

	log := log.WithField("target", "whatgoesaround")

	if !shouldCrawl() {
		log.Info("crawl is disabled - stopping")
		return
	}

	log.Info("starting crawl")

	catalog, err := getCatalog(log)
	if err != nil {
		log.WithError(err).Error("could not retrieve catalog from FTP")
		return
	}

	gzr, err := gzip.NewReader(bytes.NewReader(catalog))
	if err != nil {
		log.WithError(err).Error("could not read catalog payload")
		return
	}

	reader := csv.NewReader(gzr)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1 // allow variable number of fields
	reader.LazyQuotes = true

	// products := make([]models.ProductRecord, 0, batchSize)
	totalSaved := 0

	for i := 0; ; i++ {

		rec, err := reader.Read()
		if err == io.EOF {
			log.Info("read all records")
			break
		}
		if err != nil {
			log.WithError(err).Error("could not read record")
			return
		}

		// first and last line have different number of fields
		if len(rec) != fieldsPerRecord {

			if i == 0 {
				continue
			}

			if len(rec) > 0 {
				log.Info("found end marker")
				if rec[0] == endMarker {
					break
				}
			}

			log.WithField("field_count", len(rec)).WithField("fields", fmt.Sprintf("%+#v", rec)).Warn("unexpected CSV format")
			continue
		}

		prod, keep := createProduct(rec)
		if !keep {
			continue
		}

		err = models.GetDBv2().Create(&prod).Error
		if err != nil {
			log.WithError(err).Error("could not save rakuten products to db")
			return
		}

		totalSaved++
	}

	log.WithField("total_saved", totalSaved).Info("rakuten crawl completed")
}

func shouldCrawl() bool {

	var site models.Site
	err := models.GetDBv2().Where("id = 7").First(&site).Error
	if err != nil {
		log.WithError(err).Warn("could not determine if 'whatgoesaround' site is disabled")
		return true
	}

	return site.Active
}

func createProduct(fields []string) (models.ProductRecord, bool) {

	productURL := fmt.Sprint(fields[fieldURL])

	retailPrice, err := strconv.ParseFloat(fields[fieldRetailPrice], 64)
	if err != nil {
		log.WithError(err).WithField("target", "whatgoesaround").Warn("could not parse price")
	}
	price, _ := strconv.ParseFloat(fields[fieldPrice], 64)
	if price == 0 {
		price = retailPrice
	}

	// only save price if it's in USD
	if fields[fieldCurrency] != "USD" {
		log.WithField("currency", fields[fieldCurrency]).WithField("target", "whatgoesaround").Warn("could not parse price")
		price = 0
	}

	rec := models.ProductRecord{
		Guid:        uuid.New().String(),
		SiteId:      siteID,
		Title:       fields[fieldName],
		ProductURL:  &productURL,
		Img:         fields[fieldImage],
		UrlId:       0,
		Brand:       fields[fieldBrand],
		Price:       int64(price),
		RetailPrice: int64(retailPrice),
		Color:       fields[fieldColor],
		Size:        fields[fieldSize],
		ItemNumber:  fields[fieldID],
		Approved:    true,
	}

	category, subCategory := determineCategory(fields[fieldCategory], fields[fieldSubCategory])
	if category == "" {
		return rec, false
	}

	rec.Category = category
	rec.SubCategory = subCategory

	return rec, true
}

func determineCategory(category string, subCategory string) (string, string) {

	const (
		apparelCategory = "Apparel & Accessories"
		luggageCategory = "Luggage & Bags"

		separator = "~~"
	)

	if category != luggageCategory && category != apparelCategory {
		return "", ""
	}

	if category == luggageCategory {

		fields := strings.Split(subCategory, separator)

		sub := ""
		if len(fields) > 0 {
			sub = strings.ToLower(fields[len(fields)-1])
		}

		// return last field, lowercased
		return "bags", sub
	}

	// apparel only now
	fields := strings.Split(subCategory, separator)

	if len(fields) > 0 {

		// if it's jewelry, we're done
		if fields[0] == "Jewelry" {
			return "accessories", "Jewelry"
		}

		// decide based on last field
		field := strings.ToLower(fields[len(fields)-1])

		switch field {

		case "belts", "cufflinks", "gloves & Mittens", "headbands", "hats", "neckties", "scarves & shawls", "sunglasses":
			return "accessories", ""

		case "dresses", "jumpsuits & rompers", "outerwear", "coats & jackets", "vests", "pants", "shirts & tops", "kimonos":
			return "clothing", field
		}
	}

	return "accessories", ""
}
