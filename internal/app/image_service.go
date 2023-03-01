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
)

type ImageService interface {
	UploadImage(image *multipart.FileHeader) (int64, error)
	DownloadImage(id int64) (string, error)
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

func (i *imageService) DownloadImage(id int64) (string, error) {
	return "", nil
}

func createPath(fileName string, storage string) (string, error) {
	//Get root path name
	//cwd, _ := os.Getwd()

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
