package main

import (
	"fmt"
	"message_processor/api"
	"message_processor/config"
	"net/http"
	"time"
)

func main() {

	config := config.GetConfig()

	filesaver := api.NewLineContentSaverService(
		config.ServiceConfig.LineWebhook.AccessToken,
		config.ServiceConfig.MediaArchivePath,
	)

	msgProcessor := api.NewMessageProcessorService(
		filesaver,
		config.ServiceConfig.BotMessageCollector.Baseurl,
		&http.Client{},
		config.ServiceConfig.ApiKey,
	)

	go filesaver.SaveContentWorker()

	for {
		select {
		case <-time.After(config.ServiceConfig.GetMessageInterval):
			fmt.Println("=== Start processing messages ===")
			err := msgProcessor.Process(
				config.ServiceConfig.MaximumProcessFiles,
			)

			if err != nil {
				fmt.Println("Error processing messages:", err)
			} else {
				fmt.Println("=== Finished processing messages ===")
			}
		}
	}

}
