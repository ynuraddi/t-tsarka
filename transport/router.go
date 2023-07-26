package transport

import "github.com/labstack/echo/v4"

func (s *Server) setupRouter() {
	router := echo.New()

	router.GET("/rest/substr/find", nil)

	s.router = router
}
