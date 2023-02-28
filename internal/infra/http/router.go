package http

import (
	"github.com/BogdanStaziyev/softcery-test/config/container"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/http/validators"
	MW "github.com/labstack/echo/v4/middleware"
)

func EchoRouter(s *Server, cont container.Container) {
	e := s.Echo
	e.Use(MW.Logger())

	e.Validator = validators.NewValidator()

	_ = e.Group("image")
}
