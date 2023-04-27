package httpgin

import (
	"homework10/internal/app/adsapp"
	"homework10/internal/app/userapp"
	"homework10/internal/ports/httpgin/adsport"
	"homework10/internal/ports/httpgin/userport"
	"homework10/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(port string, ad adsapp.App, user userapp.App, log logger.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	api := handler.Group("/api/v1", Logger(log), gin.Recovery())
	{
		adsport.AppRouter(api, ad)
		userport.AppRouter(api, user)
	}

	return &http.Server{Addr: port, Handler: handler}
}
