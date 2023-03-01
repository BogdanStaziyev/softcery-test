package handlers

import (
	"github.com/BogdanStaziyev/softcery-test/internal/app"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/transport/response"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type ImageHandler struct {
	is app.ImageService
}

func NewImageHandler(imageService app.ImageService) ImageHandler {
	return ImageHandler{
		is: imageService,
	}
}

func (i *ImageHandler) Download(ctx echo.Context) error {
	image, err := ctx.FormFile("image")
	if err != nil {
		log.Println(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "The image was not uploaded. Please add an image to the field and try again.")
	}
	//Check file format
	contentType := image.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		log.Println(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "The format of the submitted file is not supported. The file should be in the format of: .png or .jpeg")
	}
	name, err := i.is.DownloadImage(image)
	if err != nil {
		log.Println(name)
		return err
	}
	return nil
}

func (i *ImageHandler) Upload(ctx echo.Context) error {
	return nil
}
