package adsport

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app/adsapp"
	"homework10/internal/entities/ads"
	"log"
	"net/http"
	"strconv"
)

func createAd(a adsapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			switch err {
			case ads.ErrInvalidAdParams:
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAdById(a adsapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adId, _ := strconv.Atoi(c.Param("ad_id"))
		log.Println(adId)
		ad, err := a.GetAdById(c, int64(adId))
		if err != nil {
			switch err {
			case adrepo.ErrInvalidAdId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAdsByTitle(a adsapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		adsArr, err := a.GetAdsByTitle(c, title)
		if err != nil {
			switch err {
			case adrepo.ErrInvalidAdTitle:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdsSuccessResponse(adsArr))
	}
}

func getAllAds(a adsapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := ads.Filters{
			Status:   ads.Status(c.DefaultQuery("status", string(ads.Published))),
			Date:     c.Query("date"),
			AuthorId: c.Query("author_id"),
		}
		adsArr, err := a.GetAll(c, filters)
		if err != nil {
			switch err {
			case ads.ErrInvalidFilters:
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		log.Println(len(adsArr))
		c.JSON(http.StatusOK, AdsSuccessResponse(adsArr))
	}
}

func changeAdStatus(a adsapp.App) gin.HandlerFunc {
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
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			case adrepo.ErrInvalidAdId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			case ads.ErrInvalidAdParams:
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func updateAd(a adsapp.App) gin.HandlerFunc {
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
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			case ads.ErrInvalidAdParams:
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case adrepo.ErrInvalidAdId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func deleteAd(a adsapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		id, _ := strconv.Atoi(c.Param("ad_id"))
		err = a.DeleteAd(c, int64(id), reqBody.UserID)
		if err != nil {
			switch err {
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			case adrepo.ErrInvalidAdId:
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			case adsapp.ErrUnableToDelete:
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdDeleteSuccessResponse())
	}
}
