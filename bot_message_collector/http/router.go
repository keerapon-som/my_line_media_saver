package http

import (
	"bot_message_collector/api"
	"bot_message_collector/config"
	"bot_message_collector/repository"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

var (
	buildtime, buildcommit, version string
)

type Handler struct {
	lineWebhook *api.LineWebhookService
	jsonArchive *repository.LineJsonfileArchive
}

func ApiAuthMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("Authorization")

	if apiKey != "Bearer "+config.GetConfig().ServiceConfig.ApiKey {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}

func NewHTTPRouter(lineWebhook *api.LineWebhookService, jsonArchive *repository.LineJsonfileArchive) *fiber.App {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Use(pprof.New())

	app.Get("/version", getVersion)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "success",
		})
	})

	h := &Handler{
		lineWebhook: lineWebhook,
		jsonArchive: jsonArchive,
	}

	app.Post("/line_chat_webhook", ApiAuthMiddleware, h.LineChatWebhook)

	app.Get("/line_chat_webhook/list_filenames", ApiAuthMiddleware, h.GetJsonArchiveLists)
	app.Post("/json_archive", ApiAuthMiddleware, h.GetJsonArchives)

	app.Delete("/json_archives", ApiAuthMiddleware, h.DeleteJsonArchives)
	app.Delete("/json_archive/all", ApiAuthMiddleware, h.DeleteAllJsonArchives)

	return app

}

func getVersion(c *fiber.Ctx) error {

	versionInfo := struct {
		BuildCommit string
		BuildTime   string
		Version     string
	}{
		BuildCommit: buildcommit,
		BuildTime:   buildtime,
		Version:     version,
	}

	return c.Status(http.StatusOK).JSON(versionInfo)
}
