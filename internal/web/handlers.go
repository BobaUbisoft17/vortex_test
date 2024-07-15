package server

import (
	"net/http"
	"vortex_test/internal/model"

	"github.com/labstack/echo/v4"
)

func (s *Server) AddClient(c echo.Context) error {
	var client model.Client
	if err := c.Bind(&client); err != nil {
		return err
	}

	client, err := s.db.AddClient(client)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return c.JSON(http.StatusCreated, client)
}

func (s *Server) UpdateClient(c echo.Context) error {
	var client model.Client
	if err := c.Bind(&client); err != nil {
		return err
	}

	client, err := s.db.UpdateClient(client)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, client)
}

func (s *Server) DeleteClient(c echo.Context) error {
	var client model.Client

	if err := c.Bind(&client); err != nil {
		return err
	}

	if err := s.db.DeleteClient(client); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Client was delete successfully")
}

func (s *Server) UpdateAlgorithmStatus(c echo.Context) error {
	var algorithm model.Algorithm

	if err := c.Bind(&algorithm); err != nil {
		return err
	}

	algorithm, err := s.db.UpdateAlgorithmStatus(algorithm)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, algorithm)
}
