package database

import (
	"github.com/BogdanStaziyev/softcery-test/internal/domain"
	"github.com/upper/db/v4"
)

const imageTable = "images"

type image struct {
	ID          int64  `db:"id,omitempty"`
	ImagePath   string `db:"image_path"`
	ContentType string `db:"content_type"`
}

type ImageRepo interface {
	// SaveImage accepts an image entity, saves it to the database, creates a unique ID, and returns it
	SaveImage(image domain.Image) (int64, error)
	// GetImage accepts a unique ID, finds the image entity associated with it, and returns it
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

	//Find one image by id
	err := i.coll.Find(db.Cond{"id": id}).One(&img)
	if err != nil {
		return domain.Image{}, err
	}
	return img.mapImageToDomain(), nil
}

// mapDomainToImages converts the received entity from the upper layer into an image entity for writing to the database
func (i *imageRepo) mapDomainToImages(img domain.Image) image {
	return image{
		ImagePath:   img.Path,
		ContentType: img.ContentType,
	}
}

// mapImageToDomain converts the retrieved entity from the database into an image entity
func (i *image) mapImageToDomain() domain.Image {
	return domain.Image{
		ID:          i.ID,
		Path:        i.ImagePath,
		ContentType: i.ContentType,
	}
}
