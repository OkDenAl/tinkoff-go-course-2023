package httpgin

import (
	"homework8/internal/app/adapp"
	"homework8/internal/app/userapp"
	"homework8/internal/ports/httpgin/adsport"
	"homework8/internal/ports/httpgin/userport"
	"homework8/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, ad adapp.App, user userapp.App) Server {
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
