package storySongHandler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mikethai/just-have-time/database"
	userHandler "github.com/mikethai/just-have-time/internal/handlers/user"
	"github.com/mikethai/just-have-time/internal/model"
	"net/http"
)

type Handler struct {
	repository Repository
	httpClient HttpClient
}

func NewHandler() *Handler {
	r := NewRepository(database.DB)
	c := NewHttpClient(&http.Client{})

	return &Handler{
		repository: r,
		httpClient: c,
	}
}

type StorySong struct {
	SongID    string `json:"song_id" validate:"required"`
	Msno      string `json:"msno" validate:"required"`
	UserImage string
	Hashtags  []string `json:"hash_tags" validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(storySong StorySong) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(storySong)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (h *Handler) GetStorySongs(c *fiber.Ctx) error {
	storysMap := make(map[string]ResponseStoty)

	storySongs, _ := h.repository.List()

	for _, storySong := range storySongs {
		var newHastags []string

		songID := storySong.SongID
		msno := storySong.Msno
		songInfo, _ := h.httpClient.GetSongInfo(songID)

		for _, hashtag := range storySong.Hashtag {
			newHastags = append(newHastags, hashtag.Name)
		}

		songAlbumImage := ""
		if len(songInfo.Album.Images) > 1 {
			songAlbumImage = songInfo.Album.Images[1].Url
		}

		newResponseStorySong := ResponseStorySong{
			SongID:         songID,
			SongName:       songInfo.Name,
			SongAlbumImage: songAlbumImage,
			Artist:         songInfo.Album.Artist.Name,
			SongHashTag:    newHastags,
			CreatedAt:      int(storySong.CreatedAt.Unix()),
		}

		entry, mapExist := storysMap[msno]
		// If the key exists
		if mapExist {
			entry.Songs = append(entry.Songs, newResponseStorySong)
			storysMap[msno] = entry
		} else {
			var songs []ResponseStorySong
			songs = append(songs, newResponseStorySong)
			storysMap[msno] = ResponseStoty{
				Msno:         msno,
				UserImage:    "https://i.kfs.io/muser/global/" + msno + "/cropresize/300x300.jpg",
				UserHashTags: []string{"Hello", "迷妹日常"},
				Songs:        songs,
			}
		}
	}

	storysSlics := getStorysAsSlice(storysMap, c.Params("msno"))

	// If no story song is present return an error
	if len(storysSlics) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No story song present", "data": nil})
	}

	// Else return story song
	return c.JSON(fiber.Map{"status": "success", "message": "Story Song Found", "data": storysSlics})
}

func (h *Handler) CreateStorySongs(c *fiber.Ctx) error {
	storySong := new(StorySong)

	if err := c.BodyParser(&storySong); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	errors := ValidateStruct(*storySong)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	isUserExist := checkUserExists(storySong.Msno)
	if !isUserExist {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "The user not exists."})
	}

	newStorySongModel := model.StorySong{
		SongID: storySong.SongID,
		Msno:   storySong.Msno,
	}

	hashTagParameter := &CreateHashTagParameter{
		storySongModel: newStorySongModel,
		Hashtags:       storySong.Hashtags,
	}

	storySongModel, err := h.repository.CreateHashTag(hashTagParameter)

	storySongParameter := &CreateParameter{
		SongID:  storySongModel.SongID,
		Msno:    storySongModel.Msno,
		Hashtag: storySongModel.Hashtag,
	}

	newStorySong, err := h.repository.Create(storySongParameter)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create story song", "data": err})
	}

	// Return the created story song
	return c.JSON(fiber.Map{"status": "success", "message": "Created Story song", "data": newStorySong})
}

func (h *Handler) GetStorySong(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Story Song Found", "data": ""})
}

func (h *Handler) UpdateStorySong(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Story Song Updated", "data": ""})
}

func checkUserExists(msno string) bool {
	newsUserHandler := userHandler.NewHandler()
	_, err := newsUserHandler.FetchUser(msno)
	if err != nil {
		return false
	}

	return true
}
