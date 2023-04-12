package httpgin

import (
	"homework8/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	gin.Default()
	s := Server{port: port, app: gin.New()}

	log := logger.InitLog()

	api := s.app.Group("/api/v1", Logger(log), gin.Recovery())
	{
		AppRouter(api, a)
	}

	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
