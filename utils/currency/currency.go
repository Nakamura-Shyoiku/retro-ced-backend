package currency

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/ulventech/retro-ced-backend/utils/client"
)

const (
	//TODO: extract into config
	apiConvertEndpoint = "https://free.currconv.com/api/v7/convert"
	apiKey             = "a17951278b9e36540ef2"

	LabelUSD = "USD"
	LabelJPY = "JPY"
)

// ConversionRate will return the conversion rate when converting 'from' currency to the 'to currency.
func ConversionRate(from string, to string) (float64, error) {

	cli, err := client.NewDefaultClient()
	if err != nil {
		return -1, errors.Wrap(err, "could not create HTTP client for converter")
	}

	conversionKey := fmt.Sprintf("%v_%v", from, to)

	params := url.Values{}
	params.Set("apiKey", apiKey)
	params.Set("compact", "ultra")
	params.Set("q", conversionKey)

	var response map[string]interface{}
	err = cli.GetJSON(apiConvertEndpoint, params, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not retrieve currency conversion rate")
	}

	rate, ok := response[conversionKey].(float64)
	if !ok {
		return -1, errors.Wrap(err, "could not locate currency conversion rate in server response")
	}

	return rate, nil
}
