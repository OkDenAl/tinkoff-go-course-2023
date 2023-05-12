package adrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework10/internal/entities/ads"
	"testing"
)

type AdsRepoTest struct {
	Name    string
	Body    any
	Expect  any
	IsError bool
}

type AddBody struct {
	Ad *ads.Ad
}

type DeleteBody struct {
	adId int64
}

type GetByIdBody struct {
	adId int64
}

type GetByTitleBody struct {
	title string
}

type UpdateAdStatusBody struct {
	adId      int64
	newStatus bool
}

type UpdateTitleBody struct {
	adId     int64
	newTitle string
	newText  string
}

func TestAdRepository(t *testing.T) {
	tests := []AdsRepoTest{
		{
			Name: "add first ad ok", Body: AddBody{Ad: &ads.Ad{Title: "title"}},
			Expect: int64(1),
		},
		{
			Name: "delete ad ok", Body: DeleteBody{adId: 0},
			Expect: nil,
		},
		{
			Name: "get ad by id ok", Body: GetByIdBody{adId: 0},
			Expect: &ads.Ad{
				Title: "test",
				ID:    0,
			},
		},
		{
			Name: "get ad by id error", Body: GetByIdBody{adId: 2},
			Expect: ErrInvalidAdId, IsError: true,
		},
		{
			Name: "get ads by title ok", Body: GetByTitleBody{title: "test"},
			Expect: []*ads.Ad{{
				Title: "test",
				ID:    0,
			}},
		},
		{
			Name: "get ads by title error", Body: GetByTitleBody{title: "kek"},
			Expect: ErrInvalidAdTitle, IsError: true,
		},
		{
			Name: "update ad status ok", Body: UpdateAdStatusBody{newStatus: true, adId: 0},
			Expect: &ads.Ad{
				Title:     "test",
				ID:        0,
				Published: true,
			},
		},
		{
			Name: "update ad status ok", Body: UpdateTitleBody{newTitle: "new title", newText: "new text", adId: 0},
			Expect: &ads.Ad{
				Title: "new title",
				ID:    0,
				Text:  "new text",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			repo := NewForTest(map[int64]*ads.Ad{
				0: {
					Title:     "test",
					ID:        0,
					Published: false,
				},
			}, 1)
			ctx := context.Background()
			switch tc.Body.(type) {
			case AddBody:
				id, _ := repo.AddAd(ctx, tc.Body.(AddBody).Ad)
				assert.Equal(t, id, tc.Expect)
			case DeleteBody:
				err := repo.DeleteAd(ctx, tc.Body.(DeleteBody).adId)
				assert.Nil(t, err)
			case GetByIdBody:
				ad, _ := repo.GetAdById(ctx, tc.Body.(GetByIdBody).adId)
				if tc.IsError {
					assert.Error(t, tc.Expect.(error))
				} else {
					assert.Equal(t, ad.Text, tc.Expect.(*ads.Ad).Text)
				}
			case GetByTitleBody:
				adsArr, _ := repo.GetAdsByTitle(ctx, tc.Body.(GetByTitleBody).title)
				if tc.IsError {
					assert.Error(t, tc.Expect.(error))
				} else {
					assert.Len(t, adsArr, 1)
					assert.Equal(t, adsArr[0].Title, tc.Expect.([]*ads.Ad)[0].Title)
				}
			case UpdateAdStatusBody:
				ad, _ := repo.UpdateAdStatus(ctx, tc.Body.(UpdateAdStatusBody).adId, tc.Body.(UpdateAdStatusBody).newStatus)
				assert.Equal(t, ad.Published, tc.Expect.(*ads.Ad).Published)
			case UpdateTitleBody:
				ad, _ := repo.UpdateAdTitleAndText(ctx, tc.Body.(UpdateTitleBody).adId, tc.Body.(UpdateTitleBody).newTitle,
					tc.Body.(UpdateTitleBody).newText)
				assert.Equal(t, ad.Text, tc.Expect.(*ads.Ad).Text)
				assert.Equal(t, ad.Title, tc.Expect.(*ads.Ad).Title)
			}

		})
	}
}
