package chartsHandler

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

func (h *Handler) GetSongCharts(c *fiber.Ctx) error {

	result, _ := h.httpClient.GetSongCharts()
	return c.JSON(fiber.Map{"status": "success", "message": "Story Song Found", "data": result.Tracks.Data[:20]})
}
