package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteAllJsonArchives(c *fiber.Ctx) error {

	err := h.jsonArchive.DeleteAllJsonFiles()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete JSON file",
			"err":   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
