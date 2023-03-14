package storage

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/domain"

	// external
	"github.com/BogdanStaziyev/softcery-test/pkg/resizer"
)

type storage struct {
	path string
}

func NewStorage(path string) *storage {
	return &storage{
		path: path,
	}
}

func (s *storage) Save(image *multipart.FileHeader, domainImage *domain.Image) error {
	//Create current path to image
	domainImage.CreatePath(image.Filename, s.path)

	src, err := image.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//Destination
	dst, err := os.Create(domainImage.Path)
	if err != nil {
		return err
	}

	//Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

// MakeVariants takes the path to an image from RabbitMQ and creates three different versions of the image
func (s *storage) MakeVariants(path string) error {
	// open image
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// Check file format
	extension := filepath.Ext(path)
	var img image.Image
	switch extension {
	case ".png":
		// decode png into image.Image
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
	default:
		// decode jpeg into image.Image
		img, err = jpeg.Decode(file)
		if err != nil {
			return err
		}
	}

	//Close file
	err = file.Close()
	if err != nil {
		return err
	}

	//Get image width size
	width := img.Bounds().Size().X

	//Create a slice with all the options for changing the image
	coefficients := []float32{75, 50, 25}
	for _, c := range coefficients {

		// resizer the width of the image
		result := float32(width) * c / 100

		m := resizer.Resize(uint(result), img)

		//split the path to the file and specify the coefficient by which the image was resized
		res := strings.Split(path, "name=")
		newPath := fmt.Sprintf("%s%s%.0f%s", res[0], "name=", c, res[1])

		// create a new file using the new path with the reduction coefficient
		file, err = os.Create(newPath)
		if err != nil {
			return err
		}

		// write new image to file
		err = jpeg.Encode(file, m, nil)
		if err != nil {
			return err
		}
		_ = file.Close()
	}
	return nil
}
