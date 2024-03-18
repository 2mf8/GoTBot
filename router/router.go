package router

import (
	"github.com/2mf8/GoTBot/api"
	"github.com/2mf8/GoTBot/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.CORS())
	r.GET("/", api.IndexApi)
	r.POST("/login", api.TokenGetApi)
	/*r.GET("/price/:service_number/:item", middleware.JWTAuthMiddlewareOrRefreshToken(), api.PriceGetItemApi)
	r.GET("/prices/:service_number/:key", api.PriceGetItemsApi)
	r.GET("/prices/:service_number", api.PriceGetItemsAllApi)
	r.POST("/price/:service_number", middleware.JWTAuthMiddlewareOrRefreshToken(), api.PriceAddAndUpdateByItemApi)
	r.POST("/price/:service_number/:item", middleware.JWTAuthMiddlewareOrRefreshToken(), api.PriceAddAndUpdateByItemApi)
	r.DELETE("/price/:service_number/:item", middleware.JWTAuthMiddlewareOrRefreshToken(), api.PriceDeleteItemApi)*/
	return r
}
