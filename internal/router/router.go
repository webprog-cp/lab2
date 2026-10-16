package router

import (
	"net/http"

	"cp_lab1/internal/app"
	"cp_lab1/internal/config"
	"cp_lab1/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(cfg *config.Config, app *app.App) (*gin.Engine, error) {
	gin.SetMode(cfg.GinMode)
	binding.EnableDecoderDisallowUnknownFields = true

	ge := gin.New()
	if err := ge.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	registerRoutes(ge, cfg, app)

	ge.Use(gin.Logger())
	ge.Use(gin.Recovery()) // Recovers from any panic and returns 500 if there is one

	return ge, nil
}

const apiKey = "demo"

func requireAPIKey(c *gin.Context) {
	if c.GetHeader("X-API-Key") != apiKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid API key"})
		return
	}
	c.Next()
}

func registerCats(ge *gin.RouterGroup, ch *handlers.CatHandler) {
	cats := ge.Group("/cats")
	cats.GET("", ch.GetAll)
	cats.GET("/:id", ch.GetByID)

	protected := cats.Group("", requireAPIKey)
	protected.POST("", ch.Create)
	protected.PUT("/:id", ch.Update)
	protected.DELETE("", ch.DeleteAll)
	protected.DELETE("/:id", ch.Delete)
}

func registerShelters(ge *gin.RouterGroup, sh *handlers.ShelterHandler) {
	shelters := ge.Group("/shelters")
	shelters.GET("", sh.GetAll)
	shelters.GET("/:id", sh.GetByID)
	shelters.GET("/:id/cats", sh.GetCats)

	protected := shelters.Group("", requireAPIKey)
	protected.POST("", sh.Create)
	protected.PUT("/:id", sh.Update)
	protected.DELETE("", sh.DeleteAll)
	protected.DELETE("/:id", sh.Delete)
}

func registerRoutes(ge *gin.Engine, cfg *config.Config, app *app.App) {
	api := ge.Group(cfg.BasePath())
	registerCats(api, app.CatHandler)
	registerShelters(api, app.ShelterHandler)
	ge.StaticFile("/openapi.yaml", "docs/openapi.yaml")
	ge.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.yaml")))
}
