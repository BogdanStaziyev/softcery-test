package domain

import (
	"fmt"
	"github.com/google/uuid"
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

func (i *Image) ReturnCurrentPath(quantity string) {
	//Split base path
	res := strings.Split(i.Path, "name=")

	//Change base path  to current version
	i.Path = fmt.Sprintf("%s%s%s%s", res[0], "name=", quantity, res[1])
}

func (i *Image) CreatePath(fileName, storage string) {
	// Create a new file name by combining the uuid and the default name. And use "name=" as a delimiter.
	newFileName := fmt.Sprintf("%sname=%s", uuid.New().String(), fileName)

	// Create file path
	path := filepath.Join(storage, newFileName)
	i.Path = filepath.FromSlash(path)
}
