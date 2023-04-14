package adsadapter

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app/adapp"
	"homework8/internal/entities/ads"
	"net/http"
	"strconv"
)

func createAd(a adapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			switch err {
			case ads.ErrInvalidAdParams:
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
		return
	}
}

func getAdById(a adapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adId, _ := strconv.Atoi(c.Param("ad_id"))
		ad, err := a.GetAdById(c, int64(adId))
		if err != nil {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
		return
	}
}

func getAdByTitle(a adapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		ad, err := a.GetAdByTitle(c, title)
		if err != nil {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
		return
	}
}

func getAllAds(a adapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		ad, err := a.GetAdByTitle(c, title)
		if err != nil {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
		return
	}
}

func changeAdStatus(a adapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		adId, _ := strconv.Atoi(c.Param("ad_id"))
		ad, err := a.ChangeAdStatus(c, int64(adId), reqBody.UserID, reqBody.Published)
		if err != nil {
			switch err {
			case ads.ErrUserCantChangeThisAd:
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
		return
	}
}

func updateAd(a adapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		adId, _ := strconv.Atoi(c.Param("ad_id"))
		ad, err := a.UpdateAd(c, int64(adId), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			switch err {
			case ads.ErrUserCantChangeThisAd:
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			case ads.ErrInvalidAdParams:
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
		return
	}
}
