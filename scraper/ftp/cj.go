package ftp

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils/bounds"
	"github.com/ulventech/retro-ced-backend/utils/ternary"
)

const (
	remoteDir = "/outgoing/productcatalog/236006/"
)

// NOTE: we currently only want to process US products
var requiredFilesRe = regexp.MustCompile(`Vestiaire_Collective.*US.*`)

// DownloadCJ will read the products from the specified SFTP server and store them in the database
func DownloadCJ() error {

	log.Info("starting CJ crawl")

	// NOTE: since the return value of the function is not checked, we sometimes both log and return error;
	// should be fixed in the future

	addr := "datatransfer.cj.com:22"
	config := &ssh.ClientConfig{
		User: "5278852",
		Auth: []ssh.AuthMethod{
			ssh.Password("r$D6iUET"),
		},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{"ssh-dss"},
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.WithError(err).Fatal("failed to dial")
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.WithError(err).Fatal("failed to create client")
	}

	// close connection when done
	defer client.Close()

	// get file list
	files, err := client.ReadDir(remoteDir)
	if err != nil {
		log.WithError(err).Fatal("could not read remote directory")
	}

	// here we try to continue despite errors
	for _, f := range files {

		if !requiredFilesRe.MatchString(f.Name()) {
			log.Infof("skipping based on name filter: %v", f.Name())
			continue
		}

		log.Infof("processing file: %s", f.Name())

		// open remote file
		file, err := client.Open(fmt.Sprintf("%s/%s", remoteDir, f.Name()))
		if err != nil {
			log.WithError(err).
				WithField("name", f.Name()).
				Warn("could not open remote file")

			continue
		}

		// get file size
		fileinfo, err := file.Stat()
		if err != nil {
			log.WithError(err).
				WithField("name", f.Name()).
				Warn("could not stat remote file")
			continue
		}

		// download file
		log.WithField("name", f.Name()).
			WithField("size", fileinfo.Size()).
			Info("downloading...")

		bufferedFile := make([]byte, fileinfo.Size())
		_, err = file.Read(bufferedFile)
		if err != nil {
			log.WithError(err).
				WithField("name", f.Name()).
				Warn("could not read remote file")
			continue
		}

		log.WithField("name", f.Name()).Info("downloaded...")

		// create Reader for the zipped payload
		zipReader, err := zip.NewReader(bytes.NewReader(bufferedFile), fileinfo.Size())
		if err != nil {
			log.WithError(err).
				WithField("name", f.Name()).
				Warn("could not create reader for remote file")
			continue
		}

		for _, zippedFile := range zipReader.File {

			filename := zippedFile.FileHeader.Name

			log.WithField("archive", f.Name()).
				WithField("name", filename).
				Info("found file in archive")

			// unpack zip file in-memory
			// TODO: we keep an awful lot in memory at the same time;
			// this can probably be streamlined
			payload, err := readZipFile(zippedFile)
			if err != nil {
				log.WithError(err).
					WithField("archive", f.Name()).
					WithField("name", filename).
					Warn("failed to unzip file")
				continue
			}

			log.WithField("archive", f.Name()).
				WithField("name", filename).
				Info("read file in archive")

			// log.Info(string(unziped))

			csvReader := csv.NewReader(bytes.NewReader(payload))
			csvReader.Comma = '\t'

			var products = make([]*models.Product, 0)

			for i := 0; ; i++ {

				record, err := csvReader.Read()
				if err != nil {
					if err == io.EOF {
						log.WithField("name", filename).
							Info("finished reading file")
						break
					}

					log.WithError(err).
						WithField("name", filename).
						Warn("could not read line from TSV")
					break
				}

				// skip file header
				if i == 0 {
					continue
				}

				// TODO: progress update on every 5k products - remove once stable enough
				if i%5000 == 0 {
					log.WithField("count", i).
						WithField("filename", filename).
						Info("processed entries from CSV")
				}

				// Filter on only womans products
				if strings.ToLower(bounds.CheckString(record, 36)) != "female" {
					continue
				}

				priceParts := strings.Split(bounds.CheckString(record, 15), " ")
				price, err := strconv.ParseFloat(ternary.String(bounds.CheckString(priceParts, 0), "0"), 64)
				if err != nil {
					log.WithError(err).
						WithField("name", filename).
						WithField("price_string", priceParts).
						WithField("line", i).
						Warn("failed to get price")
				}

				categoryParts := strings.Split(bounds.CheckString(record, 24), ">")

				product := &models.Product{
					SiteId:           14,
					Url:              bounds.CheckString(record, 7),
					Img:              bounds.CheckString(record, 9),
					Category:         strings.ToLower(strings.TrimSpace(bounds.CheckString(categoryParts, len(categoryParts)-2))),
					ItemNumber:       bounds.CheckString(record, 4),
					Brand:            strings.ToLower(bounds.CheckString(record, 25)),
					Price:            int64(price),
					Title:            bounds.CheckString(record, 5),
					Description:      bounds.CheckString(record, 6),
					Color:            strings.ToLower(bounds.CheckString(record, 35)),
					Size:             bounds.CheckString(record, 39),
					ProductCondition: strings.ToLower(bounds.CheckString(record, 29)),
					Approved:         true,
				}

				// Queue
				log.Debugf("Queued product: %+v", product)
				products = append(products, product)
			}

			log.WithField("name", filename).
				WithField("count", len(products)).
				Info("read products")

			// Save products to db
			for _, product := range products {
				// Save the url
				url := new(models.Url)
				err := url.AddUrl(uint64(product.SiteId), product.Url, product.Category)
				if err != nil {
					log.WithError(err).
						Errorf("failed to add url: %s", product.Url)
					return err
				}

				product.UrlId = int64(url.Id)

				// Save the product
				if err = product.Create(); err != nil {
					log.WithError(err).Errorf("failed to add product: %s", product.ItemNumber)
					return err
				}

				log.WithField("guid", product.Guid).Debug("saved product")
			}
		}
	}

	log.Info("completed CJ crawl")

	return nil
}
