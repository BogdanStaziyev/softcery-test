package usecase

import (
	"mime/multipart"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
)

type ImageService interface {
	// UploadImage receives a multipart.FileHeader and image entity and copy image to storage.
	UploadImage(image *multipart.FileHeader, domainImage domain.Image) (int64, error)
	// DownloadImage receives the image ID and the desired size, sends it to the database.
	//changes the received path to the desired size
	DownloadImage(id int64, quantity string) (domain.Image, error)
}

type ImageRepo interface {
	// SaveImage accepts an image entity, saves it to the database, creates a unique ID, and returns it
	SaveImage(image domain.Image) (int64, error)
	// GetImage accepts a unique ID, finds the image entity associated with it, and returns it
	GetImage(id int64) (domain.Image, error)
}

type Queue interface {
	// PublishImage accept path to image and send it to queue
	PublishImage(path string) error
	// CreateQueue create new queue
	CreateQueue() error
	// Consumer gets the image path and redirects it to create variants
	Consumer() error
}
