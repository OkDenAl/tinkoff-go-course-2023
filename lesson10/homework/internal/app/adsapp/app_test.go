package adsapp

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	adMocks "homework10/internal/app/adsapp/mocks"
	uMocks "homework10/internal/app/userapp/mocks"
	"homework10/internal/entities/ads"
	"log"
	"testing"
)

type AdServiceAddAdTestSuite struct {
	suite.Suite
	adRepo  *adMocks.Repository
	uRepo   *uMocks.Repository
	service App
}

func (suite *AdServiceAddAdTestSuite) SetupTest() {
	suite.adRepo = &adMocks.Repository{}
	suite.uRepo = &uMocks.Repository{}

	suite.adRepo.On("AddAd", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*ads.Ad")).
		Return(int64(0), nil)
}

func (suite *AdServiceAddAdTestSuite) TearDownTest() {
	log.Println("TEST DONE")
}

func (suite *AdServiceAddAdTestSuite) TestCreateAd_OK() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil)
	suite.service = NewApp(suite.adRepo, suite.uRepo)

	ad, err := suite.service.CreateAd(context.Background(), "title", "text", int64(0))
	suite.Nil(err)
	suite.Zero(ad.ID)
	suite.Equal(ad.Text, "text")
}

func (suite *AdServiceAddAdTestSuite) TestCreateAd_InvalidUserId() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, userrepo.ErrInvalidUserId)
	suite.service = NewApp(suite.adRepo, suite.uRepo)

	_, _ = suite.service.CreateAd(context.Background(), "title", "text", int64(1))
	suite.Error(userrepo.ErrInvalidUserId)
}

func (suite *AdServiceAddAdTestSuite) TestCreateAd_InvalidAdParams() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil).Maybe()
	suite.service = NewApp(suite.adRepo, suite.uRepo)

	_, _ = suite.service.CreateAd(context.Background(), "", "", int64(0))
	suite.Error(ads.ErrInvalidAdParams)
}

func TestAdAdd(t *testing.T) {
	suite.Run(t, new(AdServiceAddAdTestSuite))
}

func TestAdGetById(t *testing.T) {
	adRepo := &adMocks.Repository{}

	adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0)}, nil)
	service := NewApp(adRepo, &uMocks.Repository{})
	ad, err := service.GetAdById(context.Background(), int64(0))
	assert.Nil(t, err)
	assert.Equal(t, ad.Title, "test")
	assert.Equal(t, ad.ID, int64(0))
}

func TestAdGetByTitle(t *testing.T) {
	adRepo := &adMocks.Repository{}

	adRepo.On("GetAdsByTitle", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Return([]*ads.Ad{{Title: "test 1", Text: "test", ID: int64(0)}, {Title: "test 2", Text: "test", ID: int64(0)}}, nil)
	service := NewApp(adRepo, &uMocks.Repository{})
	all, err := service.GetAdsByTitle(context.Background(), "test")
	assert.Nil(t, err)
	assert.Len(t, all, 2)
}

func TestAdGetAll(t *testing.T) {
	adRepo := &adMocks.Repository{}

	adRepo.On("GetAll", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("ads.Filters")).
		Return([]*ads.Ad{{Title: "test2", Text: "test2", ID: int64(0)}, {Title: "test1", Text: "test1", ID: int64(0)}}, nil)

	service := NewApp(adRepo, &uMocks.Repository{})
	all, err := service.GetAll(context.Background(), ads.Filters{AuthorId: "0", Status: ads.Published})
	assert.Nil(t, err)
	assert.Len(t, all, 2)

	_, _ = service.GetAll(context.Background(), ads.Filters{AuthorId: "0"})
	assert.Error(t, ads.ErrInvalidFilters)
}

type AdServiceChangeAdTestSuite struct {
	suite.Suite
	adRepo *adMocks.Repository
	uRepo  *uMocks.Repository
}

func (suite *AdServiceChangeAdTestSuite) SetupTest() {
	suite.adRepo = &adMocks.Repository{}
	suite.uRepo = &uMocks.Repository{}

	suite.adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0), Published: false}, nil)
	suite.adRepo.On("UpdateAdStatus", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("bool")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0), Published: true}, nil)
}

func (suite *AdServiceChangeAdTestSuite) TearDownTest() {
	log.Println("TEST DONE")
}

func (suite *AdServiceChangeAdTestSuite) TestChangeAdStatus_OK() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil)
	service := NewApp(suite.adRepo, suite.uRepo)

	ad, err := service.ChangeAdStatus(context.Background(), int64(0), int64(0), true)
	suite.Nil(err)
	suite.Zero(ad.ID)
	suite.Equal(ad.Published, true)
	suite.Equal(ad.Text, "test")
}

func (suite *AdServiceChangeAdTestSuite) TestChangeAdStatus_InvalidUserId() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, userrepo.ErrInvalidUserId)
	service := NewApp(suite.adRepo, suite.uRepo)

	_, _ = service.ChangeAdStatus(context.Background(), int64(0), int64(0), true)
	suite.Error(userrepo.ErrInvalidUserId)
}

func (suite *AdServiceChangeAdTestSuite) TestChangeAdStatus_ErrUserCantChangeThisAd() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil)
	service := NewApp(suite.adRepo, suite.uRepo)

	_, _ = service.ChangeAdStatus(context.Background(), int64(0), int64(1), true)
	suite.Error(ads.ErrUserCantChangeThisAd)
}

func TestChangeAd(t *testing.T) {
	suite.Run(t, new(AdServiceChangeAdTestSuite))
}

type AdServiceUpdateAdTextTestSuite struct {
	suite.Suite
	adRepo *adMocks.Repository
	uRepo  *uMocks.Repository
}

func (suite *AdServiceUpdateAdTextTestSuite) SetupTest() {
	suite.adRepo = &adMocks.Repository{}
	suite.uRepo = &uMocks.Repository{}

	suite.adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0), Published: false}, nil)
	suite.adRepo.On("UpdateAdTitleAndText", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&ads.Ad{Title: "test new", Text: "test new", ID: int64(0), Published: true}, nil)
}

func (suite *AdServiceUpdateAdTextTestSuite) TestUpdateAdText_OK() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil)
	service := NewApp(suite.adRepo, suite.uRepo)

	ad, err := service.UpdateAd(context.Background(), int64(0), int64(0), "test new", "test new")
	suite.Nil(err)
	suite.Zero(ad.ID)
	suite.Equal(ad.Text, "test new")
	suite.Equal(ad.Title, "test new")
}

func (suite *AdServiceUpdateAdTextTestSuite) TestUpdateAdText_InvalidUserId() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, userrepo.ErrInvalidUserId)
	service := NewApp(suite.adRepo, suite.uRepo)

	_, _ = service.UpdateAd(context.Background(), int64(0), int64(0), "test new", "test new")
	suite.Error(userrepo.ErrInvalidUserId)
}

func (suite *AdServiceUpdateAdTextTestSuite) TestUpdateAdText_ErrUserCantChangeThisAd() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil)
	service := NewApp(suite.adRepo, suite.uRepo)
	_, _ = service.UpdateAd(context.Background(), int64(0), int64(1), "test new", "test new")
	suite.Error(ads.ErrUserCantChangeThisAd)
}

func (suite *AdServiceUpdateAdTextTestSuite) TestUpdateAdText_InvalidParams() {
	service := NewApp(suite.adRepo, suite.uRepo)

	_, _ = service.UpdateAd(context.Background(), int64(0), int64(0), "", "")
	suite.Error(ads.ErrInvalidAdParams)
}

func TestUpdateAdText(t *testing.T) {
	suite.Run(t, new(AdServiceUpdateAdTextTestSuite))
}

type AdServiceDeleteTestSuite struct {
	suite.Suite
	adRepo *adMocks.Repository
	uRepo  *uMocks.Repository
}

func (suite *AdServiceDeleteTestSuite) SetupTest() {
	suite.adRepo = &adMocks.Repository{}
	suite.uRepo = &uMocks.Repository{}

	suite.adRepo.On("DeleteAd", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil)
}

func (suite *AdServiceDeleteTestSuite) TestDeleteAd_OK() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, nil)
	suite.adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0), Published: false}, nil)
	service := NewApp(suite.adRepo, suite.uRepo)

	err := service.DeleteAd(context.Background(), int64(0), int64(0))
	suite.Nil(err)
}

func (suite *AdServiceDeleteTestSuite) TestDeleteAd_InvalidUserId() {
	suite.uRepo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, userrepo.ErrInvalidUserId)
	suite.adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0), Published: false}, nil)
	service := NewApp(suite.adRepo, suite.uRepo)

	_ = service.DeleteAd(context.Background(), int64(0), int64(0))
	suite.Error(userrepo.ErrInvalidUserId)
}

func (suite *AdServiceDeleteTestSuite) TestDeleteAd_ErrUnableToDelete() {
	suite.adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "test", Text: "test", ID: int64(0), Published: false}, nil)
	service := NewApp(suite.adRepo, suite.uRepo)

	_ = service.DeleteAd(context.Background(), int64(0), int64(1))
	suite.Error(ErrUnableToDelete)
}

func (suite *AdServiceDeleteTestSuite) TestDeleteAd_InvalidAdId() {
	suite.adRepo.On("GetAdById", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil, adrepo.ErrInvalidAdId)
	service := NewApp(suite.adRepo, suite.uRepo)

	_ = service.DeleteAd(context.Background(), int64(0), int64(0))
	suite.Error(adrepo.ErrInvalidAdId)
}

func TestDeleteAd(t *testing.T) {
	suite.Run(t, new(AdServiceDeleteTestSuite))
}
