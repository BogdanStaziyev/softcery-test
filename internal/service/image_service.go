package service

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/database"
	"github.com/BogdanStaziyev/softcery-test/internal/rabbit"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type ImageService interface {
	UploadImage(image *multipart.FileHeader, domainImage domain.Image) (int64, error)
	DownloadImage(id int64, quantity string) (domain.Image, error)
}

type imageService struct {
	ir      database.ImageRepo
	storage string
	mq      rabbit.Rabbit
}

func NewImageService(storage string, imageRepo database.ImageRepo, rabbit rabbit.Rabbit) ImageService {
	return &imageService{
		storage: storage,
		ir:      imageRepo,
		mq:      rabbit,
	}
}

func (i *imageService) UploadImage(image *multipart.FileHeader, domainImage domain.Image) (int64, error) {
	//Create current path to image storage
	path, err := createPath(image.Filename, i.storage)
	if err != nil {
		return 0, err
	}
	src, err := image.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	//Destination
	dst, err := os.Create(path)
	if err != nil {
		return 0, err
	}

	//Copy
	if _, err = io.Copy(dst, src); err != nil {
		return 0, err
	}
	domainImage.Path = path
	//Save image to DB
	id, err := i.ir.SaveImage(domainImage)
	if err != nil {
		return 0, err
	}

	//Send the path to RabbitMQ to create different versions
	if err = i.mq.PublishImage(path); err != nil {
		return 0, err
	}
	return id, err
}

func (i *imageService) DownloadImage(id int64, quantity string) (domain.Image, error) {
	//Getting path to default image from DB
	image, err := i.ir.GetImage(id)
	if err != nil {
		return domain.Image{}, err
	}
	switch quantity {
	case "100":
		return image, nil
	case "75":
		image.Path, err = returnCurrentPath(image.Path, "0.75")
		if err != nil {
			return domain.Image{}, err
		}
		return image, nil
	case "50":
		image.Path, err = returnCurrentPath(image.Path, "0.50")
		if err != nil {
			return domain.Image{}, err
		}
		return image, nil
	case "25":
		image.Path, err = returnCurrentPath(image.Path, "0.75")
		if err != nil {
			return domain.Image{}, err
		}
		return image, nil
	}
	return domain.Image{}, err
}

func returnCurrentPath(path string, quantity string) (string, error) {
	//Split base path
	res := strings.Split(path, "name=")

	//Change base path  to current version
	newPath := fmt.Sprintf("%s%s%s%s", res[0], "name=", quantity, res[1])

	//Check existing file
	img, err := os.Open(newPath)
	if err != nil {
		return "", err
	}
	defer img.Close()
	return newPath, nil
}

func createPath(fileName string, storage string) (string, error) {
	//Open current storage or create if not exist
	_, err := os.Open(storage)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(storage, os.ModePerm)
			if err != nil {
				return "", err
			} else {
				log.Println("Created new storage")
			}
		}
	}

	//Create a new file name by combining the uuid and the default name. And use "name=" as a delimiter.
	newFileName := fmt.Sprintf("%sname=%s", uuid.New().String(), fileName)
	path := filepath.Join(storage, newFileName)
	newFilePath := filepath.FromSlash(path)
	return newFilePath, nil
}
