package ads

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidFilters = errors.New("some field of filters is invalid")
)

type Filters struct {
	Status   Status
	Date     string
	AuthorId string
}

type Status string

const (
	Published   Status = "published"
	Unpublished Status = "unpublished"
)

func (f *Filters) ValidateFilters() error {
	if !isDateValid(f.Date) || !isStatusValid(f.Status) || !isAuthorIdValid(f.AuthorId) {
		return ErrInvalidFilters
	}
	return nil
}

func isDateValid(date string) bool {
	if date == "" {
		return true
	}
	splited := strings.Split(date, "-")
	if len(splited) != 3 {
		return false
	}

	year, errY := strconv.Atoi(splited[0])
	month, errM := strconv.Atoi(splited[1])
	day, errD := strconv.Atoi(splited[2])
	if errY != nil || errM != nil || errD != nil {
		return false
	}

	if year < 0 || year > time.Now().Year() {
		return false
	}
	if month < 0 || month > 12 {
		return false
	}
	if day < 0 || day > 31 {
		return false
	}
	return true
}

func isStatusValid(status Status) bool {
	switch status {
	case Unpublished:
	case Published:
	default:
		return false
	}
	return true
}

func isAuthorIdValid(id string) bool {
	if id == "" {
		return true
	}
	_, err := strconv.Atoi(id)
	return err == nil
}

func isAuthorIdValidIncorrect(id string) bool {
	if id == "1" {
		return false
	}
	_, err := strconv.Atoi(id)
	return err == nil
}

func isAuthorIdValid2(id string) bool {
	if id == "" {
		return true
	}
	_, err := strconv.ParseInt(id, 10, 32)
	return err == nil
}

func isAuthorIdValid3(id string) bool {
	if id == "" {
		return true
	}
	_, err := strconv.ParseInt(id, 10, 32)
	return err == nil
}
