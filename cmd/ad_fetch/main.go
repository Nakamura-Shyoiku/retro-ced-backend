package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils/env"
)

// TODO: at one point config loading should be fixed
func init() {
	env.InitEnv()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(env.GetEnv().String())
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../../config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	models.Connect()
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("invalid usage: specify a URL to process")
	}

	address := os.Args[1]

	parsed, err := url.Parse(address)
	if err != nil {
		log.Fatalf("invalid URL specified: %v", address)
	}

	if !parsed.IsAbs() {
		log.Fatalf("please specify an absolute URL")
	}

	switch parsed.Hostname() {

	case "www.tradesy.com", "tradesy.com":
		err = fetchTradesy(address)

	case "www.malleries.com", "malleries.com":
		err = fetchMalleries(address)

	default:
		log.Fatalf("unknown site: %v", parsed.Hostname())
	}

	if err != nil {
		log.Fatalf("error processing ad: %v", err)
	}
}
