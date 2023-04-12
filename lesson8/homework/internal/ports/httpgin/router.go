package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/ads", createAd(a))
	r.GET("/ads/id/:ad_id", getAdById(a))
	r.GET("/ads/title/:title", getAdByTitle(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id/text", updateAd(a))
}
