package app

import (
	"context"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, id int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adId, userId int64, newStatus bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adId, userId int64, title, text string) (*ads.Ad, error)
}

type app struct {
	repo Repository
}

func NewApp(repo Repository) App {
	return app{repo: repo}
}

func (a app) CreateAd(ctx context.Context, title, text string, userId int64) (*ads.Ad, error) {
	ad := &ads.Ad{
		Title:     title,
		Text:      text,
		AuthorID:  userId,
		Published: false,
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

func (a app) ChangeAdStatus(ctx context.Context, adId, userId int64, newStatus bool) (*ads.Ad, error) {
	ad, err := a.repo.GetAd(ctx, adId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != userId {
		return nil, ads.ErrUserCantChangeThisAd
	}
	if ad.Published == newStatus {
		return ad, nil
	}

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
	ad, err := a.repo.GetAd(ctx, adId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != userId {
		return nil, ads.ErrUserCantChangeThisAd
	}
	if ad.Text == newText && ad.Title == newTitle {
		return ad, nil
	}

	ad, err = a.repo.UpdateAdTitleAndText(ctx, adId, newTitle, newText)
	if err != nil {
		return nil, err
	}
	return ad, err
}

type Repository interface {
	AddAd(ctx context.Context, ad *ads.Ad) (int64, error)
	GetAd(ctx context.Context, adId int64) (*ads.Ad, error)
	UpdateAdStatus(ctx context.Context, adId int64, newStatus bool) (*ads.Ad, error)
	UpdateAdTitleAndText(ctx context.Context, adId int64, newTitle, newText string) (*ads.Ad, error)
}
