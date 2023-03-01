package database

import (
	"github.com/upper/db/v4"
)

const imageTable = "images"

type image struct {
	ID        int    `db:"id,omitempty"`
	ImageName string `db:"name"`
}

type ImageRepo interface {
	SaveImage(imageName string) (int64, error)
}

type imageRepo struct {
	coll db.Collection
}

func NewImageRepo(dbSession db.Session) ImageRepo {
	return &imageRepo{
		coll: dbSession.Collection(imageTable),
	}
}

func (i *imageRepo) SaveImage(imageName string) (int64, error) {
	var img image
	img.ImageName = imageName
	res, err := i.coll.Insert(&img)
	if err != nil {
		return 0, err
	}
	return res.ID().(int64), err
}
