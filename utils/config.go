package utils

import (
	"strings"

	"github.com/apex/log"
	"github.com/spf13/viper"

	"github.com/ulventech/retro-ced-backend/utils/env"
)

func InitConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(env.GetEnv().String())
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
}
