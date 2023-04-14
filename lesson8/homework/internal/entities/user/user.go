package user

import (
	"errors"
	"github.com/OkDenAl/validator"
)

var (
	ErrInvalidUserParams = errors.New("invalid user params")
)

type User struct {
	Id       int64
	Nickname string
	Email    string
	Password string
}

type ValidatorUser struct {
	NicknameMin string `validate:"min:3"`
	NicknameMax string `validate:"max:20"`
	PasswordMin string `validate:"min:5"`
	PasswordMax string `validate:"max:20"`
}

func ValidateUser(u *User) error {
	vAd := ValidatorUser{
		NicknameMin: u.Nickname,
		NicknameMax: u.Nickname,
		PasswordMin: u.Password,
		PasswordMax: u.Password,
	}
	return validator.Validate(vAd)
}
