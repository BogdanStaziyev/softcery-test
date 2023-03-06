package usecase

import (
	"io"
	"mime/multipart"
	"os"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
)

type imageService struct {
	storage string
	ir      ImageRepo
	mq      Queue
}

func NewImageService(storage string, imageRepo ImageRepo, rabbit Queue) *imageService {
	return &imageService{
		storage: storage,
		ir:      imageRepo,
		mq:      rabbit,
	}
}

func (i *imageService) UploadImage(image *multipart.FileHeader, domainImage domain.Image) (int64, error) {
	//Create current path to image storage
	err := domainImage.CreatePath(image.Filename, i.storage)
	if err != nil {
		return 0, err
	}
	src, err := image.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	//Destination
	dst, err := os.Create(domainImage.Path)
	if err != nil {
		return 0, err
	}

	//Copy
	if _, err = io.Copy(dst, src); err != nil {
		return 0, err
	}

	//Save image to PostgreSQL
	id, err := i.ir.SaveImage(domainImage)
	if err != nil {
		return 0, err
	}

	//Send the path to RabbitMQ to create different versions
	if err = i.mq.PublishImage(domainImage.Path); err != nil {
		return 0, err
	}
	return id, err
}

func (i *imageService) DownloadImage(id int64, quantity string) (domain.Image, error) {
	//Getting path to default image from PostgreSQL
	image, err := i.ir.GetImage(id)
	if err != nil {
		return domain.Image{}, err
	}

	// Find current image by quantity
	switch quantity {
	case "100":
		return image, nil
	case "75":
		err = image.ReturnCurrentPath("0.75")
		if err != nil {
			return domain.Image{}, err
		}
		return image, nil
	case "50":
		err = image.ReturnCurrentPath("0.50")
		if err != nil {
			return domain.Image{}, err
		}
		return image, nil
	case "25":
		err = image.ReturnCurrentPath("0.25")
		if err != nil {
			return domain.Image{}, err
		}
		return image, nil
	}
	return domain.Image{}, err
}
