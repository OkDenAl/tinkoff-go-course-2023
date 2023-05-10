package httpgin

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app/adsapp"
	"homework10/internal/entities/ads"
	adsServiceMock "homework10/internal/ports/httpgin/adsport/mocks"
	"homework10/internal/ports/httpgin/userport/mocks"
	"homework10/pkg/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type adData struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	AuthorID     int64  `json:"author_id"`
	CreationDate string `json:"creation_date"`
	UpdateDate   string `json:"update_date"`
	Published    bool   `json:"published"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

type deleteAdResponse struct {
	Data string `json:"data"`
}

type AdsApiTestSuite struct {
	suite.Suite
	adsService *adsServiceMock.App
	client     *http.Client
	baseURL    string
}

func (suite *AdsApiTestSuite) SetupTest() {
	suite.adsService = &adsServiceMock.App{}
	server := NewHTTPServer(":18080", suite.adsService, &mocks.App{}, logger.InitLog())
	testServer := httptest.NewServer(server.Handler)
	suite.client = testServer.Client()
	suite.baseURL = testServer.URL
}

func (suite *AdsApiTestSuite) TestCreateAd_OK() {
	suite.adsService.On("CreateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "title", Text: "text", AuthorID: 1, ID: 1}, nil)
	body := map[string]any{
		"user_id": 1,
		"title":   "title",
		"text":    "text",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/ads", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd adResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Equal(responseAd.Data.Text, "text")
	suite.Equal(responseAd.Data.ID, int64(1))
}

func (suite *AdsApiTestSuite) TestCreateAd_InvalidParams() {
	suite.adsService.On("CreateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(nil, ads.ErrInvalidAdParams)
	body := map[string]any{
		"user_id": 1,
		"title":   "",
		"text":    "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/ads", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestCreateAd_InvalidUserUd() {
	suite.adsService.On("CreateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(nil, userrepo.ErrInvalidUserId)
	body := map[string]any{
		"user_id": 1,
		"title":   "",
		"text":    "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/ads", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestCreateAd_UnexpectedError() {
	suite.adsService.On("CreateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(nil, errors.New("server error"))
	body := map[string]any{
		"user_id": 1,
		"title":   "",
		"text":    "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/ads", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *AdsApiTestSuite) TestCreateAd_NilBody() {
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/ads", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestGetAdById_OK() {
	suite.adsService.On("GetAdById", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{Title: "title", Text: "text", AuthorID: 1, ID: 10}, nil)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads/id/10", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd adResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Equal(responseAd.Data.Text, "text")
	suite.Equal(responseAd.Data.ID, int64(10))
}

func (suite *AdsApiTestSuite) TestGetAdById_InvalidAdId() {
	suite.adsService.On("GetAdById", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(nil, adrepo.ErrInvalidAdId)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads/id/10", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestGetAdById_UnexpectedError() {
	suite.adsService.On("GetAdById", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(nil, errors.New("unexpected error"))
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads/id/10", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *AdsApiTestSuite) TestGetAdsByTitle_OK() {
	suite.adsService.On("GetAdsByTitle", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string")).
		Return([]*ads.Ad{{Title: "example", Text: "text", AuthorID: 1, ID: 10}}, nil)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads/title/example", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd adsResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Len(responseAd.Data, 1)
	suite.Equal(responseAd.Data[0].Title, "example")
	suite.Equal(responseAd.Data[0].ID, int64(10))
}

func (suite *AdsApiTestSuite) TestGetAdsByTitle_InvalidAdTitle() {
	suite.adsService.On("GetAdsByTitle", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string")).
		Return(nil, adrepo.ErrInvalidAdTitle)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads/title/example", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestGetAdsByTitle_UnexpectedError() {
	suite.adsService.On("GetAdsByTitle", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string")).
		Return(nil, errors.New("unexpected error"))
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads/title/example", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *AdsApiTestSuite) TestGetAllAds_OK() {
	suite.adsService.On("GetAll", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("ads.Filters")).
		Return([]*ads.Ad{{Title: "example", Text: "text", AuthorID: 1, ID: 10},
			{Title: "example2", Text: "text2", AuthorID: 1, ID: 10}}, nil)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads?author_id=1", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd adsResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Len(responseAd.Data, 2)
	suite.Equal(responseAd.Data[0].Title, "example")
	suite.Equal(responseAd.Data[0].ID, int64(10))
}

func (suite *AdsApiTestSuite) TestGetAllAds_InvalidFilters() {
	suite.adsService.On("GetAll", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("ads.Filters")).
		Return(nil, ads.ErrInvalidFilters)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads?autdffdvfhor_id=1", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestGetAllAds_UnexpectedError() {
	suite.adsService.On("GetAll", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("ads.Filters")).
		Return(nil, errors.New("unexpected error"))
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/ads?author_id=1", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_OK() {
	suite.adsService.On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(&ads.Ad{ID: 1, AuthorID: 1, Published: true}, nil)
	body := map[string]any{
		"user_id":   1,
		"published": true,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd adResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Equal(responseAd.Data.AuthorID, int64(1))
	suite.Equal(responseAd.Data.Published, true)
	suite.Equal(responseAd.Data.ID, int64(1))
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_UserCantChangeThisAd() {
	suite.adsService.On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(nil, ads.ErrUserCantChangeThisAd)
	body := map[string]any{
		"user_id":   1,
		"published": true,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusForbidden)
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_InvalidUserId() {
	suite.adsService.On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(nil, userrepo.ErrInvalidUserId)
	body := map[string]any{
		"user_id":   1,
		"published": true,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_InvalidAdId() {
	suite.adsService.On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(nil, adrepo.ErrInvalidAdId)
	body := map[string]any{
		"user_id":   1,
		"published": true,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_InvalidAdParams() {
	suite.adsService.On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(nil, ads.ErrInvalidAdParams)
	body := map[string]any{
		"user_id":   1,
		"published": true,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_UnexpectedError() {
	suite.adsService.On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(nil, errors.New("unexpected Error"))
	body := map[string]any{
		"user_id":   1,
		"published": true,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *AdsApiTestSuite) TestChangeAdStatus_NilBody() {
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/status", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestUpdateText_OK() {
	suite.adsService.On("UpdateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&ads.Ad{ID: 1, AuthorID: 1, Text: "new text", Title: "new title"}, nil)
	body := map[string]any{
		"user_id": 1,
		"title":   "new title",
		"text":    "new text",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd adResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Equal(responseAd.Data.Text, "new text")
	suite.Equal(responseAd.Data.Title, "new title")
	suite.Equal(responseAd.Data.ID, int64(1))
}

func (suite *AdsApiTestSuite) TestUpdateAdText_UserCantChangeThisAd() {
	suite.adsService.On("UpdateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, ads.ErrUserCantChangeThisAd)
	body := map[string]any{
		"user_id": 1,
		"title":   "new title",
		"text":    "new text",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusForbidden)
}

func (suite *AdsApiTestSuite) TestUpdateAdText_InvalidUserId() {
	suite.adsService.On("UpdateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, userrepo.ErrInvalidUserId)
	body := map[string]any{
		"user_id": 1,
		"title":   "new title",
		"text":    "new text",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestUpdateAdText_InvalidAdParams() {
	suite.adsService.On("UpdateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, ads.ErrInvalidAdParams)
	body := map[string]any{
		"user_id": 1,
		"title":   "",
		"text":    "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestUpdateAdText_InvalidAdId() {
	suite.adsService.On("UpdateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, adrepo.ErrInvalidAdId)
	body := map[string]any{
		"user_id": 1,
		"title":   "",
		"text":    "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestUpdateAdText_UnexpectedError() {
	suite.adsService.On("UpdateAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, errors.New("unexpected error"))
	body := map[string]any{
		"user_id": 1,
		"title":   "",
		"text":    "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *AdsApiTestSuite) TestUpdateAdText_NilBody() {
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/ads/1/text", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestDeleteAd_OK() {
	suite.adsService.On("DeleteAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64")).
		Return(nil)
	body := map[string]any{
		"user_id": 1,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/ads/1/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var responseAd deleteAdResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &responseAd)
	suite.Equal(responseAd.Data, "ad successfully deleted")
}

func (suite *AdsApiTestSuite) TestDeleteAd_NilBody() {
	suite.adsService.On("DeleteAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64")).
		Return(nil)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/ads/1/delete", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *AdsApiTestSuite) TestDeleteAd_InvalidUserId() {
	suite.adsService.On("DeleteAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64")).
		Return(userrepo.ErrInvalidUserId)
	body := map[string]any{
		"user_id": 1,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/ads/1/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestDeleteAd_InvalidAdId() {
	suite.adsService.On("DeleteAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64")).
		Return(adrepo.ErrInvalidAdId)
	body := map[string]any{
		"user_id": 1,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/ads/1/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *AdsApiTestSuite) TestDeleteAd_UnableToDelete() {
	suite.adsService.On("DeleteAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64")).
		Return(adsapp.ErrUnableToDelete)
	body := map[string]any{
		"user_id": -1,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/ads/1/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusForbidden)
}

func (suite *AdsApiTestSuite) TestDeleteAd_UnexpectedError() {
	suite.adsService.On("DeleteAd", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("int64")).
		Return(errors.New("unexpected error"))
	body := map[string]any{
		"user_id": -1,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/ads/1/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func TestAdsApi(t *testing.T) {
	suite.Run(t, new(AdsApiTestSuite))
}
