package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	"homework8/internal/app"
	"net/http"
)

func createAd(a app.App) gin.HandlerFunc {
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

func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		id := c.GetInt64("ad_id")
		ad, err := a.ChangeAdStatus(c, id, reqBody.UserID, reqBody.Published)
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

func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		id := c.GetInt64("ad_id")
		ad, err := a.UpdateAd(c, id, reqBody.UserID, reqBody.Title, reqBody.Text)
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
