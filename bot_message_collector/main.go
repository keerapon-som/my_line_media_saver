package main

import (
	"bot_message_collector/api"
	"bot_message_collector/config"
	"bot_message_collector/repository"
	"bot_message_collector/service"
	"context"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := config.GetConfig()

	repo := repository.NewLineJsonfileArchive(config.ServiceConfig.LineWebhookArchivePath)

	line := api.NewLineWebhookService(
		config.ServiceConfig.LineWebhook.AccessToken,
		5*time.Second,
		repo,
	)

	go line.RunWorker()
	// err := service.New(pgSql, numWorker).Run(ctx)
	err := service.New(line, repo).Run(ctx)

	if err != nil {
		panic(err)
	}

}
