package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)
	_, err = client.createUser("tester1", "tester1", "tester1")
	assert.NoError(t, err)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)

	_, err = client.changeAdStatus(0, 100, true)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("tester", "tester", "tester")
	assert.NoError(t, err)
	_, err = client.createUser("tester1", "tester1", "tester1")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)

	_, err = client.updateAd(1, 100, "title", "text")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	resp, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

func TestGetAdWithIncorrectId(t *testing.T) {
	client := getTestClient()

	resp, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.getAdById(2)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetAdWithIncorrectTitle(t *testing.T) {
	client := getTestClient()

	resp, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.getAdByTitle("easy hw")
	assert.ErrorIs(t, err, ErrNotFound)
}
