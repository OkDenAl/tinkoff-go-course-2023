package httpgin

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/entities/user"
	ads1ServiceMock "homework10/internal/ports/httpgin/adsport/mocks"
	user1ServiceMock "homework10/internal/ports/httpgin/userport/mocks"
	"homework10/pkg/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ErrUnexpected = errors.New("unexpected server error")

type userData struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	Data userData `json:"data"`
}

type UserApiTestSuite struct {
	suite.Suite
	userService *user1ServiceMock.App
	client      *http.Client
	baseURL     string
}

func (suite *UserApiTestSuite) SetupTest() {
	suite.userService = &user1ServiceMock.App{}
	server := NewHTTPServer(":18080", &ads1ServiceMock.App{}, suite.userService, logger.InitLog())
	testServer := httptest.NewServer(server.Handler)
	suite.client = testServer.Client()
	suite.baseURL = testServer.URL
}

func (suite *UserApiTestSuite) TestCreateUser_Ok() {
	suite.userService.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&user.User{Nickname: "nick", Email: "email", Password: "pass", Id: 0}, nil)
	body := map[string]any{
		"nickname": "nick",
		"email":    "email",
		"password": "pass",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response userResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Data.Id, int64(0))
}

func (suite *UserApiTestSuite) TestCreateUser_NilBody() {
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *UserApiTestSuite) TestCreateUser_InvalidUserParams() {
	suite.userService.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, user.ErrInvalidUserParams)
	body := map[string]any{
		"nickname": "",
		"email":    "",
		"password": "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *UserApiTestSuite) TestCreateUser_UnexpectedError() {
	suite.userService.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, ErrUnexpected)
	body := map[string]any{
		"nickname": "",
		"email":    "",
		"password": "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *UserApiTestSuite) TestChangeNickname_Ok() {
	suite.userService.On("ChangeNickname", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(&user.User{Nickname: "new nick", Email: "email", Password: "pass", Id: 1}, nil)
	body := map[string]any{
		"nickname": "new nick",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/1/nick", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response userResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Data.Nickname, "new nick")
}

func (suite *UserApiTestSuite) TestChangeNickname_NilBody() {
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/1/nick", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *UserApiTestSuite) TestChangeNickname_InvalidUserParams() {
	suite.userService.On("ChangeNickname", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(nil, user.ErrInvalidUserParams)
	body := map[string]any{
		"nickname": "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/1/nick", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *UserApiTestSuite) TestChangeNickname_InvalidUserId() {
	suite.userService.On("ChangeNickname", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(nil, userrepo.ErrInvalidUserId)
	body := map[string]any{
		"nickname": "new nick",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/-1/nick", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *UserApiTestSuite) TestChangeNickname_UnexpectedError() {
	suite.userService.On("ChangeNickname", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(nil, ErrUnexpected)
	body := map[string]any{
		"nickname": "new nick",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/-1/nick", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *UserApiTestSuite) TestUpdatePassword_Ok() {
	suite.userService.On("UpdatePassword", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(&user.User{Nickname: "nick", Email: "email", Password: "new pass", Id: 1}, nil)
	body := map[string]any{
		"password": "new pass",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/1/password", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
}

func (suite *UserApiTestSuite) TestUpdatePassword_NilBody() {
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/1/password", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *UserApiTestSuite) TestUpdatePassword_InvalidUserParams() {
	suite.userService.On("UpdatePassword", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(nil, user.ErrInvalidUserParams)
	body := map[string]any{
		"password": "",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/1/password", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
}

func (suite *UserApiTestSuite) TestUpdatePassword_InvalidUserId() {
	suite.userService.On("UpdatePassword", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(nil, userrepo.ErrInvalidUserId)
	body := map[string]any{
		"password": "new pass",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/-1/password", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *UserApiTestSuite) TestUpdatePassword_UnexpectedError() {
	suite.userService.On("UpdatePassword", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(nil, ErrUnexpected)
	body := map[string]any{
		"password": "new pass",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, suite.baseURL+"/api/v1/user/-1/password", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *UserApiTestSuite) TestDeleteUser_Ok() {
	suite.userService.On("DeleteUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(nil)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/user/1/delete", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
}

func (suite *UserApiTestSuite) TestDeleteUser_InvalidUserId() {
	suite.userService.On("DeleteUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(userrepo.ErrInvalidUserId)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/user/-1/delete", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *UserApiTestSuite) TestDeleteUser_UnexpectedError() {
	suite.userService.On("DeleteUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(ErrUnexpected)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/user/1/delete", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func (suite *UserApiTestSuite) TestGetUser_Ok() {
	suite.userService.On("GetUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(&user.User{Nickname: "nick", Email: "email", Password: "pass", Id: 1}, nil)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/user/1/get", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response userResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Data.Id, int64(1))
}

func (suite *UserApiTestSuite) TestGetUser_InvalidUserId() {
	suite.userService.On("GetUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(nil, userrepo.ErrInvalidUserId)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/user/-1/get", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
}

func (suite *UserApiTestSuite) TestGetUser_UnexpectedError() {
	suite.userService.On("GetUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(nil, ErrUnexpected)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/user/1/get", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
}

func TestUserApi(t *testing.T) {
	suite.Run(t, new(UserApiTestSuite))
}
