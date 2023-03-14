package resizer

import (
	"github.com/nfnt/resize"
	"image"
)

func Resize(width uint, img image.Image) image.Image {
	// resizer to width using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(width, 0, img, resize.Lanczos3)
	return m
}
