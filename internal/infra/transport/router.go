package transport

import (
	"github.com/BogdanStaziyev/softcery-test/internal/app/container"
	"github.com/labstack/echo/v4"
	MW "github.com/labstack/echo/v4/middleware"
)

func EchoRouter(e *echo.Echo, cont container.Container) {
	//Options
	e.Use(MW.Logger())
	e.Use(MW.Recover())

	//Routes
	v1 := e.Group("api/v1")
	imageGroup := v1.Group("/image")

	imageGroup.POST("/upload", cont.Upload)
	imageGroup.GET("/download", cont.ImageHandler.Download)
}
