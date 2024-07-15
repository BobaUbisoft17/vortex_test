// Package of web part

package server

import (
	"vortex_test/internal/model"
	"vortex_test/pkg/logging"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	server *echo.Echo
	db     Storage
	logger *logging.Logger
}

type Storage interface {
	AddClient(model.Client) (model.Client, error)
	UpdateClient(model.Client) (model.Client, error)
	DeleteClient(model.Client) error
	UpdateAlgorithmStatus(model.Algorithm) (model.Algorithm, error)
}

func New(storage Storage, logger *logging.Logger) *Server {
	serv := &Server{
		server: echo.New(),
		db:     storage,
		logger: logger,
	}

	serv.server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	serv.server.POST("/users", serv.AddClient)
	serv.server.PUT("/users", serv.UpdateClient)
	serv.server.DELETE("/users", serv.DeleteClient)
	serv.server.PUT("/algorithmStatus", serv.UpdateAlgorithmStatus)

	return serv
}

func (s *Server) Start(address string) {
	s.server.Start(address)
}
