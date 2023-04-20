package httpgin

import (
	"homework9/internal/app/adsapp"
	"homework9/internal/app/userapp"
	"homework9/internal/ports/httpgin/adsport"
	"homework9/internal/ports/httpgin/userport"
	"homework9/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, ad adsapp.App, user userapp.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}

	log := logger.InitLog()

	api := s.app.Group("/api/v1", Logger(log), gin.Recovery())
	{
		adsport.AppRouter(api, ad)
		userport.AppRouter(api, user)
	}

	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
