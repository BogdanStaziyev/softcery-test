package usecase

import (
	"mime/multipart"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
)

type imageService struct {
	ir  ImageRepo
	mq  Queue
	str Storage
}

func NewImageService(imageRepo ImageRepo, rabbit Queue, imageStorage Storage) *imageService {
	return &imageService{
		ir:  imageRepo,
		mq:  rabbit,
		str: imageStorage,
	}
}

func (i *imageService) UploadImage(image *multipart.FileHeader, domainImage *domain.Image) (int64, error) {
	//Save image to file storage
	err := i.str.Save(image, domainImage)
	if err != nil {
		return 0, err
	}

	//Save image to PostgreSQL
	id, err := i.ir.SaveImage(*domainImage)
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

	if quantity != "100" {
		// Find current image by quantity
		image.ReturnCurrentPath(quantity)
	}

	return image, nil
}
