package utils

import (
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
	"strings"
)

func MakeVariants(path string) error {
	// open image
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	//Close file
	err = file.Close()
	if err != nil {
		return err
	}

	//Get image width size
	sizeX := img.Bounds().Size().X

	//Сreate a slice with all the options for changing the image
	size := []float32{0.75, 0.5, 0.25}
	for _, r := range size {
		result := float32(sizeX) * r

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(uint(result), 0, img, resize.Lanczos3)
		res := strings.Split(path, "name=")
		newPath := fmt.Sprintf("%s%s%.2f%s", res[0], "name=", r, res[1])
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
