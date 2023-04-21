package tests

import (
	"log"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.Equal(t, response.Data.CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, response.Data.UpdateDate, "")
	assert.False(t, response.Data.Published)
}

func TestCreateAdSync(t *testing.T) {
	client := getTestClient()
	var wg sync.WaitGroup
	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)
	wg.Add(2)
	go func() {
		defer wg.Done()
		response, err := client.createAd(0, "hello", "world")
		assert.NoError(t, err)
		assert.Zero(t, response.Data.ID)
		assert.Equal(t, response.Data.Title, "hello")
		assert.Equal(t, response.Data.Text, "world")
		assert.Equal(t, response.Data.AuthorID, int64(0))
		assert.Equal(t, response.Data.CreationDate, time.Now().UTC().Format(time.DateOnly))
		assert.Equal(t, response.Data.UpdateDate, "")
		assert.False(t, response.Data.Published)
		assert.NoError(t, err)
	}()
	runtime.Gosched()
	go func() {
		defer wg.Done()
		response, err := client.createAd(0, "hello1", "world1")
		assert.NoError(t, err)
		assert.Equal(t, response.Data.ID, int64(1))
		assert.Equal(t, response.Data.Title, "hello1")
		assert.Equal(t, response.Data.Text, "world1")
		assert.Equal(t, response.Data.AuthorID, int64(0))
		assert.Equal(t, response.Data.CreationDate, time.Now().UTC().Format(time.DateOnly))
		assert.Equal(t, response.Data.UpdateDate, "")
		assert.False(t, response.Data.Published)
		assert.NoError(t, err)
	}()
	wg.Wait()
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

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)
	_, err = client.createUser("tester1", "tester1", "tester1")
	assert.NoError(t, err)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(1, response.Data.ID, true)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, response.Data.UpdateDate, time.Now().UTC().Format(time.DateOnly))
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(1, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(1, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)
	_, err = client.createUser("tester1", "tester1", "tester1")
	assert.NoError(t, err)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(1, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, response.Data.UpdateDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestGetAdById(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hi", "tinkoff")
	assert.NoError(t, err)

	response, err := client.getAdById(1)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, response.Data.Title, "hi")
	assert.Equal(t, response.Data.Text, "tinkoff")
}

func TestGetAdByTitle(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hi", "tinkoff")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hi man", "text")
	assert.NoError(t, err)

	response, err := client.getAdByTitle("hello")
	assert.NoError(t, err)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, response.Data[0].CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, response.Data[0].Text, "world")
	assert.Equal(t, response.Data[0].ID, int64(0))

	response, err = client.getAdByTitle("hi")
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)
}

func TestChangeNickname(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)

	response, err := client.changeNickname(0, "denis")
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

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)
	_, err = client.createUser("test1", "hello1@gmail.com", "world1123")
	assert.NoError(t, err)

	_, err = client.deleteUser(0)
	assert.NoError(t, err)
	_, _ = client.deleteUser(0)
	assert.Error(t, ErrNotFound)
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)
	_, err = client.createUser("test1", "hello1@gmail.com", "world1123")
	assert.NoError(t, err)

	response, err := client.getUser(1)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Id, int64(1))
	assert.Equal(t, response.Data.Email, "hello1@gmail.com")
	assert.Equal(t, response.Data.Nickname, "test1")
	assert.Equal(t, response.Data.Password, "")
}

func TestListAds(t *testing.T) {

	client := getTestClient()

	_, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	log.Println(ads.Data)
	log.Println(len(ads.Data))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("test", "hello@gmail.com", "world123")
	assert.NoError(t, err)

	_, _ = client.createAd(0, "hello", "world")
	_, _ = client.createAd(0, "hello1", "world1")
	_, _ = client.createAd(0, "hello2", "world2")

	_, err = client.deleteAd(1, 0)
	assert.NoError(t, err)

	_, _ = client.deleteAd(1, 0)
	assert.Error(t, ErrNotFound)

	_, _ = client.deleteAd(1, 1)
	assert.Error(t, ErrForbidden)

	_, _ = client.deleteUser(0)
	assert.NoError(t, err)

	_, _ = client.deleteAd(0, 0)
	assert.Error(t, ErrNotFound)
}
