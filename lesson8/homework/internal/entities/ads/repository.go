package ads

import "context"

type Repository interface {
	AddAd(ctx context.Context, ad *Ad) (int64, error)
	GetAdById(ctx context.Context, adId int64) (*Ad, error)
	GetAdsByTitle(ctx context.Context, title string) ([]*Ad, error)
	GetAll(ctx context.Context, filters Filters) ([]*Ad, error)
	UpdateAdStatus(ctx context.Context, adId int64, newStatus bool) (*Ad, error)
	UpdateAdTitleAndText(ctx context.Context, adId int64, newTitle, newText string) (*Ad, error)
}
