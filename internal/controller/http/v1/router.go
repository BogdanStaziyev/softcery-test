package v1

import (
	// echo
	"github.com/labstack/echo/v4"
	MW "github.com/labstack/echo/v4/middleware"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/app/container"

	// external
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
)

// EchoRouter create routes using the Echo router.
func EchoRouter(e *echo.Echo, cont container.Services, l logger.Interface) {
	//Options
	e.Use(MW.Logger())
	e.Use(MW.Recover())

	//Routes
	v1 := e.Group("api/v1")
	{
		newImageHandler(v1, cont.ImageService, l)
	}
}
