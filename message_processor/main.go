package main

import (
	"fmt"
	"message_processor/api"
	"message_processor/config"
	"net/http"
	"time"
)

func main() {

	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	config := config.GetConfig()

	filesaver := api.NewLineContentSaverService(
		config.ServiceConfig.LineWebhook.AccessToken,
		config.ServiceConfig.MediaArchivePath,
	)

	msgProcessor := api.NewMessageProcessorService(
		filesaver,
		config.ServiceConfig.BotMessageCollector.Baseurl,
		&http.Client{},
	)

	go filesaver.SaveContentWorker()

	for {
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("=== Start processing messages ===")
			msgProcessor.Process()

		}
	}

}
