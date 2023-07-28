package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ynuraddi/t-tsarka/pkg/validator"
)

func (s *Server) setupRouter() {
	router := echo.New()

	router.Validator = validator.NewValidator()

	router.Use(middleware.Recover())

	router.POST("/rest/substr/find", s.findSubstring)

	router.POST("/rest/email/check", s.emailCheck)
	router.POST("/rest/iin/check", s.iinCheck)

	router.POST("/rest/counter/add/:i", s.counterAdd)
	router.POST("/rest/counter/sub/:i", s.counterSub)
	router.GET("/rest/counter/val", s.counterGet)

	s.router = router
}
