package userport

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/adapters/userrepo"
	"homework8/internal/app/userapp"
	"homework8/internal/entities/user"
	"log"
	"net/http"
	"strconv"
)

func createUser(u userapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		log.Println(reqBody)
		us, err := u.CreateUser(c, reqBody.Nickname, reqBody.Email, reqBody.Password)
		if err != nil {
			switch err {
			case user.ErrInvalidUserParams:
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(us))
	}
}

func changeNickname(u userapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeNicknameRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		id, _ := strconv.Atoi(c.Param("user_id"))
		us, err := u.ChangeNickname(c, int64(id), reqBody.Nickname)
		if err != nil {
			switch err {
			case user.ErrInvalidUserParams:
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(us))
	}
}

func updatePassword(u userapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updatePasswordRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		id, _ := strconv.Atoi(c.Param("user_id"))
		us, err := u.UpdatePassword(c, int64(id), reqBody.Password)
		if err != nil {
			switch err {
			case user.ErrInvalidUserParams:
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(us))
	}
}
