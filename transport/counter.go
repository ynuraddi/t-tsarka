package transport

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) counterAdd(c echo.Context) error {
	incParam := c.Param("i")
	inc, err := strconv.Atoi(incParam)
	if err != nil {
		s.logger.Error("failed convert param of couner route", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	if err := s.service.Counter.Add(int64(inc)); err != nil {
		s.logger.Error("failed increment counter", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{Error: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) counterSub(c echo.Context) error {
	incParam := c.Param("i")
	dec, err := strconv.Atoi(incParam)
	if err != nil {
		s.logger.Error("failed convert param of couner route", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	if err := s.service.Counter.Sub(int64(dec)); err != nil {
		s.logger.Error("failed decrement counter", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{Error: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

type counterGetResponce struct {
	Value int64 `json:"counter_value"`
}

func (s *Server) counterGet(c echo.Context) error {
	val, err := s.service.Counter.Get()
	if err != nil {
		s.logger.Error("failed get value counter", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, counterGetResponce{val})
}
