package ports

import (
	"currency/inretnal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a *service.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}
	api := s.app.Group("/api")
	AppRouter(api, *a)
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
