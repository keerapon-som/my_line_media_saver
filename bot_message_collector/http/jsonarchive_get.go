package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetJsonArchives(c *fiber.Ctx) error {

	var filenames []string

	if err := c.BodyParser(&filenames); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	data, err := h.jsonArchive.GetJsonArchives(filenames)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load data from JSON file",
			"err":   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(data)
}
