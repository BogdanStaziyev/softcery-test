package app

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ImageService interface {
	DownloadImage(image *multipart.FileHeader) (string, error)
}

type imageService struct {
	storage string
}

func NewImageService(storage string) ImageService {
	return &imageService{
		storage: storage,
	}
}

func (i imageService) DownloadImage(image *multipart.FileHeader) (string, error) {
	path, err := createPath(image.Filename, 100, i.storage)
	if err != nil {
		return "", err
	}
	src, err := image.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	//Destination
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}

	//Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dst.Name(), err
}

func createPath(fileName string, quality int, storage string) (string, error) {
	cwd, _ := os.Getwd()
	_, err := os.Open(storage)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(storage, os.ModePerm)
			if err != nil {
				return "", err
			} else {
				log.Println("Created")
			}
		}
	}
	newFileName := fmt.Sprintf("%s///%v///%s", uuid.New().String(), quality, fileName)
	path := filepath.Join(cwd, storage, newFileName)
	newFilePath := filepath.FromSlash(path)
	return newFilePath, nil
}
