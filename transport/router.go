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
	router.POST("/rest/iin/check", s.iinCheck)

	router.POST("/rest/counter/add/:i", nil)
	router.POST("/rest/counter/sub/:i", nil)
	router.POST("/rest/counter/val", nil)

	s.router = router
}
