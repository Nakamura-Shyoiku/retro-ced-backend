package whatgoesaround

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/apex/log"
	"github.com/jlaffaye/ftp"
)

const (
	whatGoesAroundID = "42946"
	ftpURL           = "aftp.linksynergy.com"
	ftpUser          = "RetroCed"
	ftpPass          = "kejUvU2H"
	ftpPort          = 21

	retryLimit = 5
)

func connect(address string) (*ftp.ServerConn, error) {

	var err error

	for i := 0; i < retryLimit; i++ {

		conn, err := ftp.Dial(address, ftp.DialWithTimeout(1*time.Minute))

		if err == nil {
			return conn, nil
		}

		log.WithField("try_no", i).WithError(err).Error("could not connect to rakuten FTP")
	}

	return nil, err
}

func ftpLogin(conn *ftp.ServerConn, user string, pass string) error {

	var err error
	for i := 0; i < retryLimit; i++ {

		err = conn.Login(user, pass)
		if err == nil {
			return nil
		}

		log.WithField("try_no", i).WithError(err).Error("could not connect to rakuten FTP")
	}

	return err
}

func getCatalog(log *log.Entry) ([]byte, error) {

	conn, err := connect(fmt.Sprintf("%v:%v", ftpURL, ftpPort))
	if err != nil {
		return nil, err
	}

	log.Info("connect ok")

	// close FTP connection
	defer conn.Quit()

	err = ftpLogin(conn, ftpUser, ftpPass)
	if err != nil {
		return nil, err
	}

	log.Info("login ok")

	entries, err := conn.List("/")
	if err != nil {
		return nil, fmt.Errorf("could not list ftp: %w", err)
	}

	log.Info("list ok")

	var catalogName string
	for _, e := range entries {

		if catalogRe.MatchString(e.Name) {
			log.WithField("name", e.Name).Info("found catalog file")
			catalogName = e.Name
			break
		}
	}

	if catalogName == "" {
		return nil, fmt.Errorf("could not locate catalog file")
	}

	r, err := conn.Retr(catalogName)
	if err != nil {
		return nil, fmt.Errorf("could not fetch catalog file: %w", err)
	}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return buf, nil
}
