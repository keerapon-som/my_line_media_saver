package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"message_processor/entities"
	"net/http"
)

type MessageProcessorService struct {
	Line_webhook_url string
	httpClient       *http.Client
	FileSaverService *FileSaverService
}

func NewMessageProcessorService(FileSaverService *FileSaverService, Line_webhook_url string, httpClient *http.Client) *MessageProcessorService {
	return &MessageProcessorService{
		Line_webhook_url: Line_webhook_url,
		httpClient:       httpClient,
		FileSaverService: FileSaverService,
	}
}

func (mc *MessageProcessorService) GetJsonArchiveLists() ([]string, error) {

	url := mc.Line_webhook_url + "/line_chat_webhook/list_filenames"

	resp, err := mc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get filenames: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get filenames, status code: %d", resp.StatusCode)
	}
	// Read the response body into a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var listsFiles []string
	err = json.Unmarshal(body, &listsFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listsFiles, nil
}

func (mc *MessageProcessorService) GetJsonArchives(filenames []string) (map[string][]entities.LineWebhook, error) {
	url := mc.Line_webhook_url + "/json_archive"

	payload := filenames
	byteString, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filenames: %w", err)
	}

	resp, err := mc.httpClient.Post(url, "application/json", bytes.NewReader(byteString))
	if err != nil {
		return nil, fmt.Errorf("failed to post filenames: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get filenames, status code: %d", resp.StatusCode)
	}

	// Read the response body into a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response map[string][]entities.LineWebhook
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

func (mc *MessageProcessorService) DeleteFiles(filenames []string) error {
	url := mc.Line_webhook_url + "/json_archives"

	// Marshal the payload
	byteString, err := json.Marshal(filenames)
	if err != nil {
		return fmt.Errorf("failed to marshal filenames: %w", err)
	}

	// Create a DELETE request with a body
	req, err := http.NewRequest("DELETE", url, bytes.NewReader(byteString))
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete files, status code: %d", resp.StatusCode)
	}

	return nil
}

func (mc *MessageProcessorService) Process() error {

	listsFiles, err := mc.GetJsonArchiveLists()
	if err != nil {
		return fmt.Errorf("failed to get JSON archive lists: %w", err)
	}

	MapFilenameChatdatas, err := mc.GetJsonArchives(listsFiles)
	if err != nil {
		return fmt.Errorf("failed to get JSON archives: %w", err)
	}

	for filename, chatDatas := range MapFilenameChatdatas {
		fmt.Println("FileName:", filename)
		for _, chatData := range chatDatas {
			for _, event := range chatData.Events {
				switch event.Message.Type {
				case "image":
					mc.FileSaverService.SendToSaveContent(event.Message.ID, event.Message.ID)
				case "audio":
					mc.FileSaverService.SendToSaveContent(event.Message.ID, event.Message.ID)
				case "text":
				case "video":
					mc.FileSaverService.SendToSaveContent(event.Message.ID, event.Message.ID)
				case "location":
				case "sticker":
				case "file":
					mc.FileSaverService.SendToSaveContent(event.Message.ID, event.Message.ID)
				}
			}
		}
		mc.DeleteFiles([]string{filename})
	}

	return nil
}
