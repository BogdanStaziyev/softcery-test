package v1

import (
	"github.com/BogdanStaziyev/softcery-test/internal/app/container"
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
	"github.com/labstack/echo/v4"
	MW "github.com/labstack/echo/v4/middleware"
)

// EchoRouter create routes using the Echo router.
func EchoRouter(e *echo.Echo, cont container.Container, l logger.Interface) {
	//Options
	e.Use(MW.Logger())
	e.Use(MW.Recover())

	//Routes
	v1 := e.Group("api/v1")
	{
		newImageHandler(v1, cont.Services.ImageService, l)
	}
}
