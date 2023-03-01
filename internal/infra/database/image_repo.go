package database

import (
	"github.com/upper/db/v4"
)

const imageTable = "images"

type image struct {
	ID        int    `db:"id,omitempty"`
	ImagePath string `db:"image_path"`
}

type ImageRepo interface {
	SaveImage(imageName string) (int64, error)
	GetFullSizePath(id int64) (string, error)
}

type imageRepo struct {
	coll db.Collection
}

func NewImageRepo(dbSession db.Session) ImageRepo {
	return &imageRepo{
		coll: dbSession.Collection(imageTable),
	}
}

func (i *imageRepo) SaveImage(path string) (int64, error) {
	var img image
	img.ImagePath = path

	//Insert to db image
	res, err := i.coll.Insert(&img)
	if err != nil {
		return 0, err
	}
	return res.ID().(int64), err
}

func (i *imageRepo) GetFullSizePath(id int64) (string, error) {
	var img image

	//Find one image
	err := i.coll.Find(db.Cond{"id": id}).One(&img)
	if err != nil {
		return "", err
	}
	return img.ImagePath, nil
}
