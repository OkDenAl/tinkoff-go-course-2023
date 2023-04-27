package userport

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/app/userapp"
)

func AppRouter(r *gin.RouterGroup, u userapp.App) {
	r.POST("/user", createUser(u))
	r.GET("/user/:user_id/get", getUser(u))
	r.PUT("/user/:user_id/nick", changeNickname(u))
	r.PUT("/user/:user_id/password", updatePassword(u))
	r.DELETE("/user/:user_id/delete", deleteUser(u))
}
