package client

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"astuart.co/goq"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// Config allows specifying the lcient configuration
type Config struct {
	Retries uint `validate:"min=1"`
	Timeout time.Duration
}

// Client provides basic HTTP functionality
type Client struct {
	client  *http.Client
	retries uint
}

// NewDefaultClient creates a client with a default configuration
func NewDefaultClient() (*Client, error) {

	defaultCfg := Config{
		Retries: 5,
		Timeout: 30 * time.Second,
	}

	return NewClient(defaultCfg)
}

// NewProxyClient creates a client with a proxy configuration
func NewProxyClient() (*Client, error) {

	defaultCfg := Config{
		Retries: 5,
		Timeout: 30 * time.Second,
	}

	// http://186.233.186.60:8080 - another proxy
	proxy := "http://20.47.108.204:8888" // need premium proxy for increasing scratching efficient
	return NewClient(defaultCfg, proxy)
}

// NewClient creates a new client based on the specified configuration
func NewClient(config Config, proxy ...string) (*Client, error) {

	err := validator.New().Struct(config)
	if err != nil {
		return nil, errors.Wrap(err, "invalid client configuration")
	}

	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Proxy connection for avoid blocking
	if len(proxy) > 0 {
		proxyUrl, err := url.Parse(proxy[0])
		if err != nil {
			return nil, errors.Wrap(err, "proxy connection failed")
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	cli := &Client{
		client:  client,
		retries: config.Retries,
	}

	return cli, nil
}

// Get performs a GET request to the specified URL with the given parameters
func (cli *Client) Get(target string, params url.Values) (*http.Response, error) {

	// parse and prepare the URL
	reqURL, err := url.Parse(target)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL format")
	}

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	//	req, err := http.NewRequest("GET", "https://api.myip.com/", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}
	req.Close = true
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	// perform the actual request
	var failures []error
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < int(cli.retries); i++ {

		// Set random User-Agent and Time interval(5s) for avoid blocking
		//req.Header.Set("User-Agent", GetUserAgent())
		time.Sleep(time.Duration(int64(rand.Intn(5))) * time.Millisecond)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36")

		res, err := cli.client.Do(req)
		if err != nil {
			failures = append(failures, err)
			continue
		}

		// we perhaps shouldn't retry on this
		if res.StatusCode != http.StatusOK {
			failures = append(failures, errors.Errorf("unexpected status code (%v)", res.Status))
			continue
		}

		return res, nil
	}

	// we failed on all retries
	messages := make([]string, 0, len(failures))
	for _, err := range failures {
		messages = append(messages, err.Error())
	}

	compositeMessage := strings.Join(messages, ", ")
	return nil, errors.Errorf("%v failures: %v", len(failures), compositeMessage)
}

// UnpackHTML will read the HTTP response body and unpack it into an HTML object
func (cli *Client) UnpackHTML(r *http.Response, v interface{}) error {

	defer r.Body.Close()

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read response body")
	}

	err = goq.Unmarshal(payload, v)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal response body")
	}

	return nil
}

// UnpackJSON will read the HTTP response and unpack it into a JSON struct
func (cli *Client) UnpackJSON(r *http.Response, v interface{}) error {

	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read response body")
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal response body")
	}

	return nil
}

// GetHTML is a shorthand for Get() and UnpackHTML()
func (cli *Client) GetHTML(target string, params url.Values, out interface{}) error {

	res, err := cli.Get(target, params)
	if err != nil {
		return err
	}

	return cli.UnpackHTML(res, out)
}

// GetJSON is a shorthand for Get() and UnpackJSON()
func (cli *Client) GetJSON(targetURL string, params url.Values, out interface{}) error {

	res, err := cli.Get(targetURL, params)
	if err != nil {
		return err
	}

	return cli.UnpackJSON(res, out)
}
