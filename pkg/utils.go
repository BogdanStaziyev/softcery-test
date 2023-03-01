package pkg

import (
	"fmt"
	"github.com/h2non/bimg"
	"log"
	"strings"
)

func MakeVariants(name string) error {
	//Read the data of the image using its path
	data, err := bimg.Read(name)
	if err != nil {
		return err
	}
	options := bimg.Options{
		Quality: 75,
	}
	image, err := bimg.NewImage(data).Process(options)
	newPath := strings.Split(name, "name=")
	log.Println(fmt.Sprintf(newPath[0] + "name=" + newPath[1]))
	err = bimg.Write(name, image)
	if err != nil {
		return err
	}
	return nil

}
