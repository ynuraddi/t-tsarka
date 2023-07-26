package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-tsarka/pkg/validator"
)

func (s *Server) setupRouter() {
	router := echo.New()

	router.Validator = validator.NewValidator()

	router.POST("/rest/substr/find", s.findSubstring)

	router.POST("/rest/email/check", s.emailCheck)

	s.router = router
}
