package utils

import (
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
	"strings"
)

// MakeVariants takes the path to an image from RabbitMQ and creates three different versions of the image
func MakeVariants(path string) error {
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
	coefficients := []float32{0.75, 0.5, 0.25}
	for _, c := range coefficients {

		// resize the width of the image
		result := float32(width) * c

		// resize to width using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(uint(result), 0, img, resize.Lanczos3)

		//split the path to the file and specify the coefficient by which the image was resized
		res := strings.Split(path, "name=")
		newPath := fmt.Sprintf("%s%s%.2f%s", res[0], "name=", c, res[1])

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
