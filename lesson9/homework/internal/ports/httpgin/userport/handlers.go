package userport

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app/userapp"
	"homework9/internal/entities/user"
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

func deleteUser(u userapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("user_id"))
		err := u.DeleteUser(c, int64(id))
		if err != nil {
			switch err {
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserDeleteSuccessResponse())
	}
}

func getUser(u userapp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("user_id"))
		u, err := u.GetUser(c, int64(id))
		if err != nil {
			switch err {
			case userrepo.ErrInvalidUserId:
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}
