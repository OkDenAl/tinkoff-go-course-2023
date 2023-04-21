package app

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
	"homework9/internal/ports/grpc/base"
	"time"
)

var (
	ErrUnableToDelete = errors.New("unable to delete an ad created by another user")
)

type AdService struct {
	repo     ads.Repository
	userRepo user.Repository
	base.UnimplementedAdServiceServer
}

func NewAdService(repo ads.Repository, userRepo user.Repository) *AdService {
	return &AdService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (a *AdService) CreateAd(ctx context.Context, req *base.CreateAdRequest) (*base.AdResponse, error) {
	_, err := a.userRepo.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	t := time.Now().UTC()
	ad := &ads.Ad{
		Title:        req.Title,
		Text:         req.Text,
		AuthorID:     req.UserId,
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
	return &base.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationDate,
		UpdateDate:   ad.UpdateDate,
	}, nil
}

func (a *AdService) ChangeAdStatus(ctx context.Context, req *base.ChangeAdStatusRequest) (*base.AdResponse, error) {
	_, err := a.userRepo.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	ad, err := a.repo.GetAdById(ctx, req.AdId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != req.UserId {
		return nil, ads.ErrUserCantChangeThisAd
	}
	if ad.Published == req.Published {
		return &base.AdResponse{
			Id:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		}, nil
	}

	t := time.Now().UTC()
	ad.UpdateDate = t.Format(time.DateOnly)

	ad, err = a.repo.UpdateAdStatus(ctx, req.AdId, req.Published)
	if err != nil {
		return nil, err
	}
	return &base.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationDate,
		UpdateDate:   ad.UpdateDate,
	}, nil
}

func (a *AdService) UpdateAd(ctx context.Context, req *base.UpdateAdRequest) (*base.AdResponse, error) {
	err := ads.ValidateAd(&ads.Ad{Title: req.Title, Text: req.Text})
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	_, err = a.userRepo.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	ad, err := a.repo.GetAdById(ctx, req.AdId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != req.UserId {
		return nil, ads.ErrUserCantChangeThisAd
	}
	if ad.Text == req.Text && ad.Title == req.Title {
		return &base.AdResponse{
			Id:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		}, nil
	}

	t := time.Now().UTC()
	ad.UpdateDate = t.Format(time.DateOnly)

	ad, err = a.repo.UpdateAdTitleAndText(ctx, req.AdId, req.Title, req.Text)
	if err != nil {
		return nil, err
	}
	return &base.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationDate,
		UpdateDate:   ad.UpdateDate,
	}, nil
}

func (a *AdService) GetAdById(ctx context.Context, req *base.GetAdByIdRequest) (*base.AdResponse, error) {
	ad, err := a.repo.GetAdById(ctx, req.AdId)
	if err != nil {
		return nil, err
	}
	return &base.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationDate,
		UpdateDate:   ad.UpdateDate,
	}, nil
}

func (a *AdService) GetAdByTitle(ctx context.Context, req *base.GetAdByTitleRequest) (*base.ListAdResponse, error) {
	adsArr, err := a.repo.GetAdsByTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	response := make([]*base.AdResponse, len(adsArr))
	for i, ad := range adsArr {
		response[i] = &base.AdResponse{
			Id:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		}
	}
	return &base.ListAdResponse{List: response}, nil
}

func (a *AdService) ListAds(ctx context.Context, f *base.Filters) (*base.ListAdResponse, error) {
	filters := ads.Filters{
		Status:   ads.Status(f.Status),
		Date:     f.Date,
		AuthorId: f.AuthorId,
	}
	err := filters.ValidateFilters()
	if err != nil {
		return nil, err
	}
	adsArr, err := a.repo.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	response := make([]*base.AdResponse, len(adsArr))
	for i, ad := range adsArr {
		response[i] = &base.AdResponse{
			Id:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		}
	}
	return &base.ListAdResponse{List: response}, nil
}

func (a *AdService) DeleteAd(ctx context.Context, req *base.DeleteAdRequest) (*empty.Empty, error) {
	ad, err := a.repo.GetAdById(ctx, req.AdId)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != req.AuthorId {
		return nil, ErrUnableToDelete
	}
	_, err = a.userRepo.GetUser(ctx, req.AuthorId)
	if err != nil {
		return nil, userrepo.ErrInvalidUserId
	}
	return &empty.Empty{}, a.repo.DeleteAd(ctx, req.AdId)
}
