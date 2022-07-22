package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ulventech/retro-ced-backend/controllers"
	"github.com/ulventech/retro-ced-backend/middlewares"
	"github.com/ulventech/retro-ced-backend/scraper/impact"
	"github.com/ulventech/retro-ced-backend/scraper/rakuten/whatgoesaround"
)

// GetEngine returns configured router
func GetEngine() *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Total-Count"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
	}))
	r.Use(middlewares.Trace())

	// Register api end-points
	v1 := r.Group("/v1")
	{
		mailchimp := new(controllers.Mailchimp)
		v1.POST("/mailchimp", mailchimp.AddSubscriber)

		products := new(controllers.Product)
		v1.GET("/product/categories/:category/:offset", products.GetProductsByCategory)
		v1.GET("/product/featured", products.GetProductsByFeatured)
		v1.GET("/product/brands/:brand/:offset", products.GetProductsByBrand)
		v1.GET("/product/brand", products.GetBrands)
		v1.GET("/product/search/:search/:offset", products.SearchProducts)
		v1.GET("/product/favourites", products.GetFavourites)
		v1.POST("/product/favourites/:product_id", products.AddFavourite)
		v1.DELETE("product/favourites/:product_id", products.RemoveFavourite)
		v1.GET("/product/is_favorited/:product_id", products.IsFavourited)
		v1.GET("/product/by-id/:id", products.GetProduct)

		trending := new(controllers.Trending)
		v1.GET("/trending", trending.GetPosts)

		support := new(controllers.Support)
		v1.POST("/support/contact-us", support.ContactUs)
		v1.POST("/support/partner-application", support.PartnerApplication)

		auth := new(controllers.Auth)
		v1.PUT("/auth", auth.Login)
		v1.POST("/auth", auth.Register)
		v1.GET("/auth", auth.ValidateToken)
		v1.GET("/auth/fb/callback", auth.FbCallback)
		v1.GET("/auth/user", auth.GetUserDetails)
		v1.POST("/auth/password-reset", auth.ResetPassword)
		v1.POST("/auth/set-new-password", auth.SetNewPassword)

		analytics := new(controllers.Analytics)
		v1.GET("/link", analytics.TrackLinkClick)

		admin := new(controllers.Admin)
		v1.GET("/admin/products", products.GetProducts)
		v1.GET("/admin/products-count", products.GetProductsCount)
		v1.PUT("/admin/products/setfeatured", admin.SetFeatured)
		v1.GET("/admin/sites", admin.GetSites)
		v1.PUT("/admin/site", admin.UpdateSite)
		v1.GET("/admin/click_tracking/summary", admin.ClickTrackingSummary)
		v1.GET("/admin/products/:status/", admin.GetProductsByApprovedStatus)
		v1.PUT("/admin/product", admin.UpdateApprovedStatus)
		v1.GET("/admin/product/:product_id", admin.GetProductByID)
		v1.PUT("/admin/product/:product_id", admin.UpdateProductById)
		v1.DELETE("/admin/product/:product_id", admin.DeleteProductById)
		v1.DELETE("/admin/product", admin.DeleteProductsById)
		v1.GET("/admin/users", admin.GetUsers)
		v1.PUT("/admin/user/:userID/acl", admin.ChangeUserPermission)
		v1.GET("/admin/partner", admin.GetPartnerTrackingSummary)
		v1.GET("/admin/partner/product", admin.GetProductTrackingSummary)
		v1.DELETE("/admin/site/:siteID", admin.DeleteSite)
		search := controllers.Search{}
		v1.POST("/search", search.SearchProducts)

		// used to get possible filter values for multiple-select filters
		v1.POST("/search/brand", search.SearchBrand)
		v1.POST("/search/size", search.SearchSize)
		v1.POST("/search/sub_category", search.SearchSubcategory)
		v1.POST("/search/color", search.SearchColor)
	}

	// The scraper is started by a call from a appengine CRON job
	scraper := new(controllers.Scraper)
	r.GET("/scraper/start", scraper.RunScrapers)
	r.GET("/scraper/sftp/start", scraper.RunSFTPScraper)
	r.GET("/scraper/impact/start", impact.CrawlTheRealReal)
	r.GET("/scraper/whatgoesaround/start", whatgoesaround.Crawl)

	// remove stale products and URLs
	r.GET("/prune", controllers.Cleanup)

	r.GET("/_ah/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"message": "404 page not found",
			},
		})
	})

	return r
}
