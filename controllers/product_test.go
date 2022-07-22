package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/repos"
	"github.com/ulventech/retro-ced-backend/utils"
)

func TestProduct_GetProducts(t *testing.T) {
	token := utils.NewToken("111", "foo@bar.com", "foo", "foo", "foo", "", 0)

	viper.SetConfigFile("../config/test.yaml")
	//viper.AutomaticEnv()
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//viper.SetConfigName(env.GetEnv().String())
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("../config/")
	//viper.AddConfigPath("../../../config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	models.Connect()

	router := gin.Default()

	repo, err := repos.NewClickhouseProductsRepo("clickhouse://130.211.211.73:9000/retroced?user=default&password=default")
	assert.NoError(t, err)

	p := Product{ProductsRepo: repo}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", token)
	router.GET("/", p.GetProducts)
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, r)
	var prods []models.Product
	err = json.Unmarshal(writer.Body.Bytes(), &prods)
	assert.NoError(t, err)
	assert.NotEmpty(t, prods)

	assert.True(t, func() bool {
		for _, pitem := range prods {
			if strings.Contains(pitem.ProductURL, "shop.rebagg.com") {
				return true
			}
		}
		return false
	}())
}
