package storySongHandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mikethai/just-have-time/database"
	"github.com/mikethai/just-have-time/internal/model"
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
	SongID   string
	Msno     int64
	Hashtags []string
}

func (h *Handler) GetStorySongs(c *fiber.Ctx) error {
	storysMap := make(map[int64]ResponseStoty)
	var songs []ResponseStorySong

	storySongs, _ := h.repository.List()

	for _, storySong := range storySongs {
		var newHastags []string

		songID := storySong.SongID
		msno := storySong.Msno
		songInfo, _ := h.httpClient.GetSongInfo(songID)

		for _, hashtag := range storySong.Hashtag {
			newHastags = append(newHastags, hashtag.Name)
		}

		newResponseStorySong := ResponseStorySong{
			SongID:      songID,
			SongName:    songInfo.Name,
			Artist:      songInfo.Album.Artist.Name,
			SongHashTag: newHastags,
			CreatedAt:   int(storySong.CreatedAt.Unix()),
		}
		songs = append(songs, newResponseStorySong)

		entry, mapExist := storysMap[msno]
		// If the key exists
		if mapExist {
			entry.Songs = append(entry.Songs, newResponseStorySong)
			storysMap[msno] = entry
		} else {
			storysMap[msno] = ResponseStoty{
				Msno:         msno,
				UserHashTags: []string{"Hello", "迷妹日常"},
				Songs:        songs,
			}
		}
	}

	storysSlics := getStorysAsSlice(storysMap)

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
