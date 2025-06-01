package http

import (
	"bot_message_collector/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) LineChatWebhook(c *fiber.Ctx) error {

	var receivedData entities.LineWebhook

	if err := c.BodyParser(&receivedData); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	h.lineWebhook.SendToChan(receivedData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
