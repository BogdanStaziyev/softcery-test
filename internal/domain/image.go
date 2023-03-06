package domain

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
)

// Image type should contain:
// ID which is auto-generated in the database as an int64 type
// Path to the saved image
// ContentType for proper image retrieval
type Image struct {
	ID          int64  `db:"id,omitempty"`
	Path        string `db:"image_path"`
	ContentType string `db:"content_type"`
}

func (i *Image) ReturnCurrentPath(quantity string) error {
	//Split base path
	res := strings.Split(i.Path, "name=")

	//Change base path  to current version
	newPath := fmt.Sprintf("%s%s%s%s", res[0], "name=", quantity, res[1])

	//Check existing file
	img, err := os.Open(newPath)
	if err != nil {
		return err
	}
	defer img.Close()
	i.Path = newPath
	return nil
}

func (i *Image) CreatePath(fileName, storage string) error {
	// Create a new file name by combining the uuid and the default name. And use "name=" as a delimiter.
	newFileName := fmt.Sprintf("%sname=%s", uuid.New().String(), fileName)

	// Create file path
	path := filepath.Join(storage, newFileName)
	newFilePath := filepath.FromSlash(path)
	i.Path = newFilePath
	return nil
}
