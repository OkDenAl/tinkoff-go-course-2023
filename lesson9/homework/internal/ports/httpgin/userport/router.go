package userport

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/app/userapp"
)

func AppRouter(r *gin.RouterGroup, u userapp.App) {
	r.POST("/user", createUser(u))
	r.PUT("/user/:user_id/nick", changeNickname(u))
	r.PUT("/user/:user_id/password", updatePassword(u))
}
