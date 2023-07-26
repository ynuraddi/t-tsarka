package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type findSubstrRequest struct {
	Input string `json:"input" validate:"required"`
}

type findSubstrResponce struct {
	Output string `json:"output"`
}

func (s *Server) findSubstring(ctx echo.Context) error {
	var req findSubstrRequest
	if err := ctx.Bind(&req); err != nil {
		s.logger.Error("failed bind request", err)
		return ctx.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}
	s.logger.Debug("success binding request - FindSubstring()")

	result := s.service.Substr.FindSubstr(req.Input)

	return ctx.JSON(http.StatusOK, findSubstrResponce{result})
}
