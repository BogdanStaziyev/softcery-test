package app

import (
	"fmt"
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
	UploadImage(image *multipart.FileHeader) (int64, error)
	DownloadImage(id int64, quantity string) (string, error)
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

func (i *imageService) UploadImage(image *multipart.FileHeader) (int64, error) {
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

	//Save image to DB
	id, err := i.ir.SaveImage(path)
	if err != nil {
		return 0, err
	}

	//Send the path to RabbitMQ to create different versions
	if err = i.mq.PublishImage(path); err != nil {
		return 0, err
	}
	return id, err
}

func (i *imageService) DownloadImage(id int64, quantity string) (string, error) {
	//Getting path to default image from DB
	path, err := i.ir.GetFullSizePath(id)
	if err != nil {
		return "", err
	}
	switch quantity {
	case "100":
		return path, err
	case "75":
		return returnCurrentPath(path, "0.75")
	case "50":
		return returnCurrentPath(path, "0.50")
	case "25":
		return returnCurrentPath(path, "0.25")
	}
	return "", nil
}

func returnCurrentPath(path string, quantity string) (string, error) {
	//Split base path
	res := strings.Split(path, "name=")

	//Change base path  to current version
	newPath := fmt.Sprintf("%s%s%s%s", res[0], "name=", quantity, res[1])

	//Check existing file
	_, err := os.Open(newPath)
	if err != nil {
		return "", err
	}
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
