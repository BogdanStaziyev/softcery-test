package controllers

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/http/response"
	"github.com/BogdanStaziyev/softcery-test/internal/service"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	is service.ImageService
}

func NewImageHandler(imageService service.ImageService) ImageHandler {
	return ImageHandler{
		is: imageService,
	}
}

// Upload uploading a new image, we get the image, check the format, and send it to the service layer
func (i *ImageHandler) Upload(ctx echo.Context) error {

	//Create new image entity
	var domainImage domain.Image

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
	domainImage.ContentType = contentType
	//Upload image to storage and write to DB
	imageID, err := i.is.UploadImage(image, domainImage)
	if err != nil {
		log.Println(err)
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return response.MessageResponse(ctx, http.StatusOK, fmt.Sprintf("Image successful upload, id: %d", imageID))
}

// Download retrieves an image if it exists using query parameters for ID and quantity.
// Quantity should be 75, 50 or 25 percentages
func (i *ImageHandler) Download(ctx echo.Context) error {
	//Get image id from query params
	id := ctx.QueryParams().Get("id")
	if id == "" {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "ID field is empty")
	}

	//Get quantity from query params
	quantity := ctx.QueryParam("quantity")

	//Check quantity should be one of 100, 75, 50, 25 by default quantity = 100
	if quantity == "" || quantity == "100" {
		quantity = "100"
	} else if quantity != "75" && quantity != "50" && quantity != "25" {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Error quantity should to be one of 100%, 75%, 50%, 25%")
	}
	imageID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Error ID should be integer")
	}

	//Get current path to image
	image, err := i.is.DownloadImage(imageID, quantity)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	//Return correct image to download
	err = ctx.File(image.Path)
	if err != nil {
		return err
	}
	return response.Response(ctx, http.StatusOK, fmt.Sprintf("Content-Type %s", image.ContentType))
}
