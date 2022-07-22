package main

import (
	"os"
	"runtime"

	"cloud.google.com/go/profiler"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/routes"
	"github.com/ulventech/retro-ced-backend/utils"
	"github.com/ulventech/retro-ced-backend/utils/env"
	"github.com/ulventech/retro-ced-backend/workers"
)

func init() {
	// Startup raven
	raven.SetDSN("https://9056e04ebc8d46d2bc7ab5172e48011c:f778cd8cdcf44aa6a2245b0b745d6b61@sentry.io/534615")

	// Make sure to init configs first as models depend on config
	log.SetHandler(text.New(os.Stderr))
	env.InitEnv()
	utils.InitConfig()

	// Setup google profiler
	if env.IsProd() {
		pc := profiler.Config{
			ProjectID: viper.GetString("app.project"),
			Service:   viper.GetString("app.service"),
		}
		if err := profiler.Start(pc); err != nil {
			log.WithError(err).Error("failed to start google profiler")
		}
	}

	// connect to db
	models.Connect()

	// Start workers
	workers.StartDispatcher(runtime.NumCPU())

	// Configure gin run mode
	if env.IsTest() {
		gin.SetMode(gin.TestMode)
	}
	if env.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {

	log.Info("Server start")
	log.Infof("System running in [%s] mode", env.GetEnv())

	// Start HTTP server
	routes.GetEngine().Run(":" + viper.GetString("app.port"))
}
