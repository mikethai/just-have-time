package storySongHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mikethai/just-have-time/config"
	"github.com/mikethai/just-have-time/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() ([]model.StorySong, error)
	Create(param *CreateParameter) (*model.StorySong, error)
	CreateHashTag(param *CreateHashTagParameter) (*model.StorySong, error)
	GetSongInfo(songID string) (*SongInfo, error)
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

type GetSongInfoParameter struct {
	SongID string
}

type SongInfo struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Duration int       `json:"duration"`
	Isrc     string    `json:"isrc"`
	Url      string    `json:"url"`
	Album    AlbumInfo `json:"album"`
}

type AlbumInfo struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Url         string      `json:"url"`
	ReleaseDate string      `json:"release_date"`
	Images      []ImageInfo `json:"images"`
	Artist      ArtistInfo  `json:"artist"`
}

type ImageInfo struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type ArtistInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Images []ImageInfo
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

func (r *repository) GetSongInfo(songID string) (*SongInfo, error) {
	url := "https://api.kkbox.com/v1.1/tracks/" + songID + "?territory=TW"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	bearerToken := config.Config("KKBOX_OPENAPI_BEARER_TOKEN")
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var songInfo SongInfo
	json.Unmarshal(reqBody, &songInfo)

	return &songInfo, nil
}
