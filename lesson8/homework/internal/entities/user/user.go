package user

import "github.com/OkDenAl/validator"

type User struct {
	Id       int64
	Nickname string
	Email    string
	Password string
}

type ValidatorAd struct {
	NicknameMin string `validate:"min:3"`
	NicknameMax string `validate:"max:20"`
	PasswordMin string `validate:"min:5"`
	PasswordMax string `validate:"max:20"`
}

func ValidateAd(u *User) error {
	vAd := ValidatorAd{
		NicknameMin: u.Nickname,
		NicknameMax: u.Nickname,
		PasswordMin: u.Password,
		PasswordMax: u.Password,
	}
	return validator.Validate(vAd)
}
