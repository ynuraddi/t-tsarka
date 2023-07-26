package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type findSubstrRequest struct {
	Input string `json:"input" validate:"required,alphanum"`
}

type findSubstrResponce struct {
	Output string `json:"output"`
}

func (s *Server) findSubstring(c echo.Context) error {
	var req findSubstrRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("failed bind request", err)
		return c.JSON(http.StatusUnprocessableEntity, errorResponce{err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		s.logger.Error("bad request", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	s.logger.Debug("success binding request - FindSubstring()")

	result := s.service.Substr.FindSubstr(req.Input)

	return c.JSON(http.StatusOK, findSubstrResponce{result})
}
