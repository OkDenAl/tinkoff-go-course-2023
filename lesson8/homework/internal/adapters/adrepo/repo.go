package adrepo

import (
	"context"
	"errors"
	"homework8/internal/entities/ads"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrInvalidAdId    = errors.New("cant find this id in map")
	ErrInvalidAdTitle = errors.New("cant find this title in map")
)

type repository struct {
	mu             *sync.Mutex
	adDataById     map[int64]*ads.Ad
	curIdGenerator int64
}

func New() ads.Repository {
	var mu sync.Mutex
	return &repository{adDataById: make(map[int64]*ads.Ad), curIdGenerator: 0, mu: &mu}
}

func (r *repository) AddAd(ctx context.Context, ad *ads.Ad) (int64, error) {
	//log.Println(ad.Title)
	r.mu.Lock()
	//log.Println(r.curIdGenerator, ad.Title)
	ad.ID = r.curIdGenerator
	r.adDataById[r.curIdGenerator] = ad
	r.curIdGenerator++
	//log.Println(r.curIdGenerator, ad.Title)
	r.mu.Unlock()
	return ad.ID, nil
}

func (r *repository) GetAdById(ctx context.Context, adId int64) (*ads.Ad, error) {
	r.mu.Lock()
	ad, ok := r.adDataById[adId]
	r.mu.Unlock()
	if !ok {
		return nil, ErrInvalidAdId
	}
	return ad, nil
}

func (r *repository) GetAdsByTitle(ctx context.Context, title string) ([]*ads.Ad, error) {
	resp := make([]*ads.Ad, 0)
	for key := range r.adDataById {
		r.mu.Lock()
		ad := r.adDataById[key]
		r.mu.Unlock()
		if strings.HasPrefix(ad.Title, title) {
			resp = append(resp, ad)
		}
	}
	if len(resp) == 0 {
		return nil, ErrInvalidAdTitle
	}
	return resp, nil
}

func (r *repository) GetAll(ctx context.Context, filters ads.Filters) ([]*ads.Ad, error) {
	resp := make([]*ads.Ad, 0)
	for key := range r.adDataById {
		r.mu.Lock()
		val := r.adDataById[key]
		r.mu.Unlock()
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
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adDataById[adId].Published = newStatus
	return r.adDataById[adId], nil
}

func (r *repository) UpdateAdTitleAndText(ctx context.Context, adId int64, newTitle, newText string) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adDataById[adId].Text = newText
	r.adDataById[adId].Title = newTitle
	return r.adDataById[adId], nil
}
