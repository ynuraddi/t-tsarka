package transport

import "github.com/labstack/echo/v4"

func (s *Server) setupRouter() {
	router := echo.New()

	router.POST("/rest/substr/find", s.findSubstring)

	s.router = router
}
