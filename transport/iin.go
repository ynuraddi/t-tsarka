package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type checkIINRequest struct {
	Input string `json:"input" validate:"required"`
}

type checkIINResponce struct {
	Output []string `json:"IINs"`
}

func (s *Server) iinCheck(c echo.Context) error {
	var req checkIINRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("failed bind request", err)
		return c.JSON(http.StatusUnprocessableEntity, errorResponce{err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		s.logger.Error("bad request", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	s.logger.Debug("success binding request - Check()")

	result := s.service.IIN.Check(req.Input)
	if result == nil {
		return c.JSON(http.StatusNotFound, checkEmailResponce{[]string{}})
	}

	return c.JSON(http.StatusOK, checkIINResponce{result})
}
