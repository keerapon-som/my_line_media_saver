package api

import (
	"bot_message_collector/entities"
	"bot_message_collector/repository"
	"fmt"
	"time"
)

type LineWebhookService struct {
	accessToken                string
	lineWebhookChannelStreamer chan entities.LineWebhook
	intervalSaveListData       time.Duration
	ListData                   []entities.LineWebhook
	jsonFileArchive            *repository.LineJsonfileArchive
}

func NewLineWebhookService(accessToken string, intervalSaveListData time.Duration, jsonFileArchive *repository.LineJsonfileArchive) *LineWebhookService {
	return &LineWebhookService{
		accessToken:                accessToken,
		lineWebhookChannelStreamer: make(chan entities.LineWebhook),
		ListData:                   []entities.LineWebhook{},
		intervalSaveListData:       intervalSaveListData,
		jsonFileArchive:            jsonFileArchive,
	}
}

func (l *LineWebhookService) SendToChan(data entities.LineWebhook) {
	l.lineWebhookChannelStreamer <- data
}

func (l *LineWebhookService) RunWorker() {
	fmt.Println("Starting LineWebhookService worker...")
	for {
		select {
		case msg := <-l.lineWebhookChannelStreamer:
			fmt.Println("Received message:", msg)

			l.ListData = append(l.ListData, msg)

		case <-time.After(l.intervalSaveListData):

			if len(l.ListData) > 0 {
				fmt.Println("Saving data to JSON file...")

				err := l.jsonFileArchive.SaveDataToJsonFile(l.ListData)
				if err != nil {
					fmt.Println("Error saving data to JSON file:", err)
				}

				l.ListData = []entities.LineWebhook{}
			}
		}
	}
}
