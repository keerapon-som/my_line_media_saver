package http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteJsonArchivesLower(c *fiber.Ctx) error {
	var timestamp int64
	if err := c.BodyParser(&timestamp); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"err":   err.Error(),
		})
	}
	if timestamp <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Timestamp parameter is required and must be greater than 0",
		})
	}

	deletedList, err := h.jsonArchive.DeleteListTimestampLowerThan(timestamp)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to delete JSON files",
			"err":         err.Error(),
			"deletedList": deletedList,
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"message":     fmt.Sprintf("Files with timestamp lower than %d deleted successfully", timestamp),
		"deletedList": deletedList,
	})
}

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
