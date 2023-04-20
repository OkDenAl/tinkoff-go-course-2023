package adsapp

import (
	"context"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
	"time"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, id int64) (*ads.Ad, error)
	GetAdById(ctx context.Context, id int64) (*ads.Ad, error)
	GetAdsByTitle(ctx context.Context, title string) ([]*ads.Ad, error)
	GetAll(ctx context.Context, filters ads.Filters) ([]*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adId, userId int64, newStatus bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adId, userId int64, title, text string) (*ads.Ad, error)
}

type app struct {
	repo     ads.Repository
	userRepo user.Repository
}

func NewApp(repo ads.Repository, userRepo user.Repository) App {
	return app{repo: repo, userRepo: userRepo}
}

func (a app) CreateAd(ctx context.Context, title, text string, userId int64) (*ads.Ad, error) {
	_, err := a.userRepo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	t := time.Now().UTC()
	ad := &ads.Ad{
		Title:        title,
		Text:         text,
		AuthorID:     userId,
		Published:    false,
		CreationDate: t.Format(time.DateOnly),
	}
	err = ads.ValidateAd(ad)
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

func (a app) GetAdsByTitle(ctx context.Context, title string) ([]*ads.Ad, error) {
	return a.repo.GetAdsByTitle(ctx, title)
}

func (a app) GetAll(ctx context.Context, filters ads.Filters) ([]*ads.Ad, error) {
	err := filters.ValidateFilters()
	if err != nil {
		return nil, err
	}
	return a.repo.GetAll(ctx, filters)
}

func (a app) ChangeAdStatus(ctx context.Context, adId, userId int64, newStatus bool) (*ads.Ad, error) {
	_, err := a.userRepo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
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

	t := time.Now().UTC()
	ad.UpdateDate = t.Format(time.DateOnly)

	return a.repo.UpdateAdStatus(ctx, adId, newStatus)
}

func (a app) UpdateAd(ctx context.Context, adId, userId int64, newTitle, newText string) (*ads.Ad, error) {
	err := ads.ValidateAd(&ads.Ad{Title: newTitle, Text: newText})
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	_, err = a.userRepo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
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

	t := time.Now().UTC()
	ad.UpdateDate = t.Format(time.DateOnly)

	return a.repo.UpdateAdTitleAndText(ctx, adId, newTitle, newText)
}
