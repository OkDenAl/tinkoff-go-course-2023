package adrepo

import (
	"context"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type repository struct {
	data           map[int64]*ads.Ad
	curIdGenerator int64
}

func New() app.Repository {
	return &repository{data: make(map[int64]*ads.Ad), curIdGenerator: 0}
}

func (r *repository) AddAd(ctx context.Context, ad *ads.Ad) (int64, error) {
	ad.ID = r.curIdGenerator
	r.data[r.curIdGenerator] = ad
	r.curIdGenerator++
	return ad.ID, nil
}

func (r *repository) GetAd(ctx context.Context, adId int64) (*ads.Ad, error) {
	return r.data[adId], nil
}

func (r *repository) UpdateAdStatus(ctx context.Context, adId int64, newStatus bool) (*ads.Ad, error) {
	r.data[adId].Published = newStatus
	return r.data[adId], nil
}

func (r *repository) UpdateAdTitleAndText(ctx context.Context, adId int64, newTitle, newText string) (*ads.Ad, error) {
	r.data[adId].Text = newText
	r.data[adId].Title = newTitle
	return r.data[adId], nil
}
