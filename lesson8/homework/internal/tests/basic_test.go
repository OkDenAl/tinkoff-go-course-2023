package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.Equal(t, response.Data.CreationDate, time.Now().Format(time.DateOnly))
	assert.Equal(t, response.Data.UpdateDate, "")
	assert.False(t, response.Data.Published)
}

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.Id)
	assert.Equal(t, response.Data.Email, "hello@gmail.com")
	assert.Equal(t, response.Data.Nickname, "test")
	assert.Equal(t, response.Data.Password, "")
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().Format(time.DateOnly))
	assert.Equal(t, response.Data.UpdateDate, time.Now().Format(time.DateOnly))
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(123, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(123, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(123, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().Format(time.DateOnly))
	assert.Equal(t, response.Data.UpdateDate, time.Now().Format(time.DateOnly))
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestGetAdById(t *testing.T) {
	client := getTestClient()

	_, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(1, "hi", "tinkoff")
	assert.NoError(t, err)

	response, err := client.getAdById(1)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().Format(time.DateOnly))
	assert.Equal(t, response.Data.Title, "hi")
	assert.Equal(t, response.Data.Text, "tinkoff")
}

func TestGetAdByTitle(t *testing.T) {
	client := getTestClient()

	_, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(1, "hi", "tinkoff")
	assert.NoError(t, err)

	response, err := client.getAdByTitle("hi")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().Format(time.DateOnly))
	assert.Equal(t, response.Data.Text, "tinkoff")
	assert.Equal(t, response.Data.ID, int64(1))
}

func TestChangeNickname(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)

	response, err = client.changeNickname(0, "denis")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "denis")

	response, err = client.changeNickname(0, "denis")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "denis")
}

func TestUpdatePassword(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)

	_, err = client.updatePassword(response.Data.Id, "denis")
	assert.NoError(t, err)
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}
