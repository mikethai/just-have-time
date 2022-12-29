package storySongHandler

import (
	"github.com/mikethai/just-have-time/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() ([]model.StorySong, error)
	Create(param *CreateParameter) (*model.StorySong, error)
	CreateHashTag(param *CreateHashTagParameter) (*model.StorySong, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

type CreateParameter struct {
	SongID  string
	Msno    int64
	Hashtag []model.Hashtag
}

type CreateHashTagParameter struct {
	storySongModel model.StorySong
	Hashtags       []string
}

func (r *repository) List() ([]model.StorySong, error) {

	var storySong []model.StorySong

	// find all storySong in the database
	r.db.Preload("Hashtag").Find(&storySong)

	return storySong, nil
}

func (r *repository) Create(param *CreateParameter) (*model.StorySong, error) {

	storySong := model.StorySong{
		SongID:  param.SongID,
		Msno:    param.Msno,
		Hashtag: param.Hashtag,
	}

	err := r.db.Create(&storySong).Error
	if err != nil {
		return nil, err
	}

	return &storySong, nil
}

func (r *repository) CreateHashTag(param *CreateHashTagParameter) (*model.StorySong, error) {

	for _, hashTagName := range param.Hashtags {
		var hashTag model.Hashtag
		if err := r.db.Where("name = ?", hashTagName).First(&hashTag).Error; err != nil {
			hashTag = model.Hashtag{Name: hashTagName}
			r.db.Create(&hashTag)
		}
		r.db.Model(&param.storySongModel).Association("Hashtag").Append(&hashTag)
	}

	return &param.storySongModel, nil
}
