package searchHandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	httpClient HttpClient
}

func NewHandler() *Handler {
	c := NewHttpClient(&http.Client{})

	return &Handler{httpClient: c}
}

func (h *Handler) GetSongByKeyword(c *fiber.Ctx) error {

	keyWordResult, _ := h.httpClient.GetSearchSong(c.Params("keyWord"))
	return c.JSON(fiber.Map{"status": "success", "message": "Story Song Found", "data": keyWordResult.Tracks.Data})
}
