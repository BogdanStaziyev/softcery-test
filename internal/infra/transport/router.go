package transport

import (
	"github.com/BogdanStaziyev/softcery-test/config/container"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/transport/validators"
	MW "github.com/labstack/echo/v4/middleware"
)

func EchoRouter(s *Server, cont container.Container) {
	e := s.Echo
	e.Use(MW.Logger())
	e.Use(MW.Recover())

	e.Validator = validators.NewValidator()

	v1 := e.Group("api/v1")
	imageGroup := v1.Group("/image")

	imageGroup.POST("/download", cont.Download)
	imageGroup.GET("/upload/:quantity", cont.ImageHandler.Upload)
}
