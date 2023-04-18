package userport

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/entities/user"
)

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type changeNicknameRequest struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type updatePasswordRequest struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
}

func UserSuccessResponse(u *user.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			Id:       u.Id,
			Email:    u.Email,
			Nickname: u.Nickname,
		},
		"error": nil,
	}
}

func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
