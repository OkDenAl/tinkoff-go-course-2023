package ads

import (
	"errors"
	"github.com/OkDenAl/validator"
)

var (
	ErrUserCantChangeThisAd = errors.New("the user is trying to change an ad created by another user")
	ErrInvalidAdParams      = errors.New("invalid ad params")
)

type Ad struct {
	ID           int64
	Title        string
	Text         string
	AuthorID     int64
	CreationDate string
	UpdateDate   string
	Published    bool
}

type ValidatorAd struct {
	TitleMin string `validate:"min:1"`
	TitleMax string `validate:"max:100"`
	TextMin  string `validate:"min:1"`
	TextMax  string `validate:"max:500"`
}

func ValidateAd(ad *Ad) error {
	vAd := ValidatorAd{
		TitleMin: ad.Title,
		TitleMax: ad.Title,
		TextMin:  ad.Text,
		TextMax:  ad.Text,
	}
	return validator.Validate(vAd)
}
