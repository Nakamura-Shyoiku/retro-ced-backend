package malleries

import (
	"github.com/pkg/errors"
)

// Fetch will retrieve and parse a single product page
func (m *Crawler) Fetch(address string) (interface{}, error) {

	res, err := m.cli.Get(address, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve page (url: %v)", address)
	}

	var page rawAdPage
	err = m.cli.UnpackHTML(res, &page)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unpack data (url: %v)", address)
	}

	return page.adPage(), nil
}
