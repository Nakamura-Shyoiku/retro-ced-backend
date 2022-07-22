package client

import (
	"fmt"
	"math/rand"
	"time"
)

var ffVersions = []float32{
	58.0,
	57.0,
	56.0,
	52.0,
	48.0,
	40.0,
	35.0,
}

var chromeVersions = []string{
	"65.0.3325.146",
	"64.0.3282.0",
	"41.0.2228.0",
	"40.0.2214.93",
	"37.0.2062.124",
}

var osStrings = []string{
	"Macintosh; Intel Mac OS X 10_10",
	"Windows NT 10.0",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",
	"X11; Linux x86_64",
}

var useragentGens = []func() string{
	genFirefoxUA, // Get User-Agent for Firefox
	genChromeUA,  // Get User-Agent for Chrome
}

func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

// Get random User-Agent (Http Header Request) for avoid blocking
func GetUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return useragentGens[rand.Intn(len(useragentGens))]()
}
