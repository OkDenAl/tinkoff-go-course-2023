package adrepo

import (
	"context"
	"errors"
	"homework8/internal/entities/ads"
	"strconv"
)

var (
	ErrInvalidAdId    = errors.New("cant find this id in map")
	ErrInvalidAdTitle = errors.New("cant find this title in map")
)

type repository struct {
	adDataById     map[int64]*ads.Ad
	adDataByTitle  map[string]*ads.Ad
	curIdGenerator int64
}

func New() ads.Repository {
	return &repository{adDataById: make(map[int64]*ads.Ad), adDataByTitle: make(map[string]*ads.Ad), curIdGenerator: 0}
}

func (r *repository) AddAd(ctx context.Context, ad *ads.Ad) (int64, error) {
	ad.ID = r.curIdGenerator

	r.adDataById[r.curIdGenerator] = ad
	r.adDataByTitle[ad.Title] = ad

	r.curIdGenerator++
	return ad.ID, nil
}

func (r *repository) GetAdById(ctx context.Context, adId int64) (*ads.Ad, error) {
	ad, ok := r.adDataById[adId]
	if !ok {
		return nil, ErrInvalidAdId
	}
	return ad, nil
}

func (r *repository) GetAdByTitle(ctx context.Context, title string) (*ads.Ad, error) {
	ad, ok := r.adDataByTitle[title]
	if !ok {
		return nil, ErrInvalidAdTitle
	}
	return ad, nil
}

func (r *repository) GetAll(ctx context.Context, filters ads.Filters) ([]*ads.Ad, error) {
	resp := make([]*ads.Ad, 0)
	for _, val := range r.adDataById {
		if filters.Date != "" && val.CreationDate != filters.Date {
			continue
		}
		if filters.AuthorId != "" && strconv.FormatInt(val.AuthorID, 10) != filters.AuthorId {
			continue
		}
		if filters.Status == ads.Published && val.Published {
			resp = append(resp, val)
		}
		if filters.Status == ads.Unpublished && !val.Published {
			resp = append(resp, val)
		}
	}
	return resp, nil
}

func (r *repository) UpdateAdStatus(ctx context.Context, adId int64, newStatus bool) (*ads.Ad, error) {
	r.adDataById[adId].Published = newStatus
	return r.adDataById[adId], nil
}

func (r *repository) UpdateAdTitleAndText(ctx context.Context, adId int64, newTitle, newText string) (*ads.Ad, error) {
	r.adDataById[adId].Text = newText
	r.adDataById[adId].Title = newTitle
	return r.adDataById[adId], nil
}
