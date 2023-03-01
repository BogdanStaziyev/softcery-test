package handlers

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/internal/app"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/transport/response"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	is app.ImageService
}

func NewImageHandler(imageService app.ImageService) ImageHandler {
	return ImageHandler{
		is: imageService,
	}
}

func (i *ImageHandler) Upload(ctx echo.Context) error {

	//Get FileHeader the multipart form file
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

	//Upload image to storage and write to DB
	imageID, err := i.is.UploadImage(image)
	if err != nil {
		log.Println(err)
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return response.MessageResponse(ctx, http.StatusOK, fmt.Sprintf("Image successful upload, id: %d", imageID))
}

func (i *ImageHandler) Download(ctx echo.Context) error {
	//Get image id from query params
	id := ctx.QueryParams().Get("id")
	if id == "" {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "ID field is empty")
	}

	//Get quantity from query params
	quantity := ctx.QueryParam("quantity")

	//Check quantity should be one of 100, 75, 50, 25 by default quantity = 100
	if quantity == "" {
		quantity = "100"
	} else if quantity != "75" && quantity != "50" && quantity != "25" {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Error quantity should to be one of 75%, 50%, 25%")
	}
	imageID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Error ID should be integer")
	}

	//Get current path to image
	path, err := i.is.DownloadImage(imageID, quantity)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return response.Response(ctx, http.StatusOK, path)
}
