package http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteJsonArchives(c *fiber.Ctx) error {
	var filenames []string
	if err := c.BodyParser(&filenames); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"err":   err.Error(),
		})
	}
	if filenames == nil || len(filenames) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Filename parameter is required",
		})
	}

	err := h.jsonArchive.DeleteJsonFile(filenames)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete JSON file",
			"err":   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("File %s deleted successfully", filenames),
	})
}
