package database

import (
	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/domain"

	// external
	"github.com/BogdanStaziyev/softcery-test/pkg/database"
)

const imageTable = "images"

type imageRepo struct {
	database.Database
}

func NewImageRepo(dbSession database.Database) *imageRepo {
	return &imageRepo{
		dbSession,
	}
}

func (i *imageRepo) SaveImage(image domain.Image) (int64, error) {

	//Insert to db image
	res, err := i.Collection(imageTable).Insert(&image)
	if err != nil {
		return 0, err
	}
	return res.ID().(int64), err
}

func (i *imageRepo) GetImage(id int64) (domain.Image, error) {
	var img domain.Image

	//Find one image by id
	err := i.Collection(imageTable).Find("id", id).One(&img)
	if err != nil {
		return domain.Image{}, err
	}
	return img, nil
}
