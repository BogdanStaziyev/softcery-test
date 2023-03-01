package database

import (
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
	"github.com/upper/db/v4"
)

const imageTable = "images"

type image struct {
	ID          int    `db:"id,omitempty"`
	ImagePath   string `db:"image_path"`
	ContentType string `db:"content_type"`
}

type ImageRepo interface {
	SaveImage(imageName domain.Image) (int64, error)
	GetImage(id int64) (domain.Image, error)
}

type imageRepo struct {
	coll db.Collection
}

func NewImageRepo(dbSession db.Session) ImageRepo {
	return &imageRepo{
		coll: dbSession.Collection(imageTable),
	}
}

func (i *imageRepo) SaveImage(image domain.Image) (int64, error) {
	img := i.mapDomainToImages(image)

	//Insert to db image
	res, err := i.coll.Insert(&img)
	if err != nil {
		return 0, err
	}
	return res.ID().(int64), err
}

func (i *imageRepo) GetImage(id int64) (domain.Image, error) {
	var img image

	//Find one image
	err := i.coll.Find(db.Cond{"id": id}).One(&img)
	if err != nil {
		return domain.Image{}, err
	}
	return img.mapImageToDomain(), nil
}

func (i *imageRepo) mapDomainToImages(img domain.Image) image {
	return image{
		ImagePath:   img.Path,
		ContentType: img.ContentType,
	}
}

func (i *image) mapImageToDomain() domain.Image {
	return domain.Image{
		ID:          i.ID,
		Path:        i.ImagePath,
		ContentType: i.ContentType,
	}
}
