package adsport

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/entities/ads"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	AuthorID     int64  `json:"author_id"`
	CreationDate string `json:"creation_date"`
	UpdateDate   string `json:"update_date"`
	Published    bool   `json:"published"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		},
		"error": nil,
	}
}

func AdsSuccessResponse(ads []*ads.Ad) *gin.H {
	resp := make([]adResponse, len(ads))
	for i, ad := range ads {
		resp[i] = adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		}
	}
	return &gin.H{
		"data":  resp,
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
