package database

import (
	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/domain"

	// external
	"github.com/BogdanStaziyev/softcery-test/pkg/database"
)

const imageTable = "images"

type ImageRepo interface {
	// SaveImage accepts an image entity, saves it to the database, creates a unique ID, and returns it
	SaveImage(image domain.Image) (int64, error)
	// GetImage accepts a unique ID, finds the image entity associated with it, and returns it
	GetImage(id int64) (domain.Image, error)
}

type imageRepo struct {
	coll database.PostgreSQL
}

func NewImageRepo(dbSession *database.PostgreSQL) ImageRepo {
	return &imageRepo{
		coll: *dbSession,
	}
}

func (i *imageRepo) SaveImage(image domain.Image) (int64, error) {

	//Insert to db image
	res, err := i.coll.Collection(imageTable).Insert(&image)
	if err != nil {
		return 0, err
	}
	return res.ID().(int64), err
}

func (i *imageRepo) GetImage(id int64) (domain.Image, error) {
	var img domain.Image

	//Find one image by id
	err := i.coll.Collection(imageTable).Find("id", id).One(&img)
	if err != nil {
		return domain.Image{}, err
	}
	return img, nil
}
