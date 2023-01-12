package storySongHandler

import (
	"time"

	"github.com/mikethai/just-have-time/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() ([]model.StorySong, error)
	Create(param *CreateParameter) (*model.StorySong, error)
	Delete(ID int) error
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
	Msno    string
	Hashtag []model.Hashtag
}

type CreateHashTagParameter struct {
	storySongModel model.StorySong
	Hashtags       []string
}

func (r *repository) List() ([]model.StorySong, error) {

	var storySong []model.StorySong

	// find all storySong in the database
	currentTime := time.Now()
	daysAgo := currentTime.Add(-time.Hour * 24)

	r.db.Model(&storySong).Preload("Hashtag").Joins("User").Where("created_at > ?", daysAgo.Format(time.RFC3339)).Order("created_at desc").Find(&storySong)

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
		} else {
			// Increment Counter
			if err := r.db.Model(&hashTag).Update("count", hashTag.Count+1).Error; err != nil {
				return nil, err
			}
		}

		r.db.Model(&param.storySongModel).Association("Hashtag").Append(&hashTag)
	}

	return &param.storySongModel, nil
}

func (r *repository) Delete(ID int) error {
	r.db.Delete(&model.StorySong{}, ID)
	return nil
}
