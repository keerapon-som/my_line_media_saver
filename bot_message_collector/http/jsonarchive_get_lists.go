package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetJsonArchiveLists(c *fiber.Ctx) error {
	filenames, err := h.jsonArchive.GetListFilenames()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve filenames",
		})

	}
	return c.Status(http.StatusOK).JSON(filenames)
}
