package v1

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/internal/controller/http/response"
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
	"github.com/BogdanStaziyev/softcery-test/internal/usecase/service"
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type imageHandler struct {
	is service.ImageService
	l  logger.Interface
}

func newImageHandler(handler *echo.Group, imageService service.ImageService, l logger.Interface) {
	r := &imageHandler{
		is: imageService,
		l:  l,
	}

	h := handler.Group("/image")
	{
		h.POST("/upload", r.Upload)
		h.GET("/download", r.Download)
	}
}

// Upload uploading a new image, we get the image, check the format, and send it to the service layer
func (i *imageHandler) Upload(ctx echo.Context) error {

	//Create new image entity
	var domainImage domain.Image

	//Get FileHeader the multipart form file
	image, err := ctx.FormFile("image")
	if err != nil {
		i.l.Error(err, "http - v1 - Upload")
		return response.ErrorResponse(ctx, http.StatusBadRequest, "The image was not uploaded. Please add an image to the field and try again.")
	}

	//Check file format
	contentType := image.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		i.l.Error(err, "http - v1 - Upload")
		return response.ErrorResponse(ctx, http.StatusBadRequest, "The format of the submitted file is not supported. The file should be in the format of: .png or .jpeg")
	}
	domainImage.ContentType = contentType
	//Upload image to storage and write to DB
	imageID, err := i.is.UploadImage(image, domainImage)
	if err != nil {
		i.l.Error(err, "http - v1 - Upload")
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return response.MessageResponse(ctx, http.StatusOK, fmt.Sprintf("Image successful upload, id: %d", imageID))
}

// Download retrieves an image if it exists using query parameters for ID and quantity.
// Quantity should be 75, 50 or 25 percentages
func (i *imageHandler) Download(ctx echo.Context) error {
	//Get image id from query params
	id := ctx.QueryParams().Get("id")
	if id == "" {
		i.l.Info("ID field is empty", "http - v1 - Download")
		return response.ErrorResponse(ctx, http.StatusBadRequest, "ID field is empty")
	}

	//Get quantity from query params
	quantity := ctx.QueryParam("quantity")

	//Check quantity should be one of 100, 75, 50, 25 by default quantity = 100
	if quantity == "" || quantity == "100" {
		quantity = "100"
	} else if quantity != "75" && quantity != "50" && quantity != "25" {
		i.l.Info("Error quantity should to be one of 100, 75, 50, 25", "http - v1 - Download")
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Error quantity should to be one of 100%, 75%, 50%, 25%")
	}
	imageID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		i.l.Error(err, "http - v1 - Download")
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Error ID should be integer")
	}

	//Get current path to image
	image, err := i.is.DownloadImage(imageID, quantity)
	if err != nil {
		i.l.Error(err, "http - v1 - Download")
		return response.ErrorResponse(ctx, http.StatusInternalServerError, "no such file")
	}

	//Return correct image to download
	err = ctx.File(image.Path)
	if err != nil {
		i.l.Error(err, "http - v1 - Download")
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return response.Response(ctx, http.StatusOK, fmt.Sprintf("Content-Type %s", image.ContentType))
}
