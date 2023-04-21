package adsport

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/app/adsapp"
)

func AppRouter(r *gin.RouterGroup, a adsapp.App) {
	r.POST("/ads", createAd(a))
	r.GET("/ads", getAllAds(a))
	r.GET("/ads/id/:ad_id", getAdById(a))
	r.GET("/ads/all", getAllAds(a))
	r.GET("/ads/title/:title", getAdsByTitle(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id/text", updateAd(a))
	r.DELETE("/ads/:ad_id/delete", deleteAd(a))
}
