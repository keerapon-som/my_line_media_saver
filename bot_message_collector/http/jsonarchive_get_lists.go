package http

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetJsonArchiveLists(c *fiber.Ctx) error {

	moreThanTimestampInt := c.QueryInt("more_than_timestamp", 0)

	log.Print("GetJsonArchiveLists called timestamp : ", moreThanTimestampInt)

	filenames, err := h.jsonArchive.GetListFilenames(moreThanTimestampInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve filenames",
		})

	}
	return c.Status(http.StatusOK).JSON(filenames)
}
