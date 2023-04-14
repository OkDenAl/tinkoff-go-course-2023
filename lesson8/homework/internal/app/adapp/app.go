package adapp

import (
	"context"
	"homework8/internal/entities/ads"
	"time"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, id int64) (*ads.Ad, error)
	GetAdById(ctx context.Context, id int64) (*ads.Ad, error)
	GetAdByTitle(ctx context.Context, title string) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adId, userId int64, newStatus bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adId, userId int64, title, text string) (*ads.Ad, error)
}

type app struct {
	repo ads.Repository
}

func NewApp(repo ads.Repository) App {
	return app{repo: repo}
}

func (a app) CreateAd(ctx context.Context, title, text string, userId int64) (*ads.Ad, error) {
	t := time.Now()
	ad := &ads.Ad{
		Title:        title,
		Text:         text,
		AuthorID:     userId,
		Published:    false,
		CreationDate: t.Format(time.DateOnly),
	}
	err := ads.ValidateAd(ad)
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	id, err := a.repo.AddAd(ctx, ad)
	if err != nil {
		return nil, err
	}
	ad.ID = id
	return ad, nil
}

func (a app) GetAdById(ctx context.Context, id int64) (*ads.Ad, error) {
	return a.repo.GetAdById(ctx, id)
}

func (a app) GetAdByTitle(ctx context.Context, title string) (*ads.Ad, error) {
	return a.repo.GetAdByTitle(ctx, title)
}

func (a app) ChangeAdStatus(ctx context.Context, adId, userId int64, newStatus bool) (*ads.Ad, error) {
	ad, err := a.repo.GetAdById(ctx, adId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != userId {
		return nil, ads.ErrUserCantChangeThisAd
	}
	if ad.Published == newStatus {
		return ad, nil
	}

	t := time.Now()
	ad.UpdateDate = t.Format(time.DateOnly)

	ad, err = a.repo.UpdateAdStatus(ctx, adId, newStatus)
	if err != nil {
		return nil, err
	}
	return ad, err
}

func (a app) UpdateAd(ctx context.Context, adId, userId int64, newTitle, newText string) (*ads.Ad, error) {
	err := ads.ValidateAd(&ads.Ad{Title: newTitle, Text: newText})
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	ad, err := a.repo.GetAdById(ctx, adId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != userId {
		return nil, ads.ErrUserCantChangeThisAd
	}
	if ad.Text == newText && ad.Title == newTitle {
		return ad, nil
	}

	t := time.Now()
	ad.UpdateDate = t.Format(time.DateOnly)

	ad, err = a.repo.UpdateAdTitleAndText(ctx, adId, newTitle, newText)
	if err != nil {
		return nil, err
	}
	return ad, err
}
