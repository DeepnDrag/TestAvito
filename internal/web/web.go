package web

import (
	"TestAvito/internal/config"
	"TestAvito/internal/storage"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/exp/slog"
	"io"
)

type Server struct {
	app     *echo.Echo
	URL     string
	logger  *slog.Logger
	Storage *storage.Storage
	JWT     config.JWT
}

func New(srvCfg config.Server, Jwt config.JWT, logger *slog.Logger, storage *storage.Storage) (*Server, error) {
	e := echo.New()
	server := Server{
		app:     e,
		URL:     srvCfg.Url,
		logger:  logger,
		Storage: storage,
		JWT:     Jwt,
	}
	e.HideBanner = true
	e.Logger.SetOutput(io.Discard)

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	m := NewMiddleware(Jwt, logger)

	server.RegisterHandlers(m)

	return &server, nil
}

func (s *Server) Serve() error {
	s.logger.Info("HTTP server started", slog.String("url", s.URL))

	return fmt.Errorf("server error: %w", s.app.Start(s.URL))
}
