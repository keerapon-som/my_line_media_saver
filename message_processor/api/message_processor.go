package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"message_processor/entities"
	"net/http"
	"os"
)

type Records struct {
	latestSavedTimestamp int64
}

type MessageProcessorService struct {
	Line_webhook_url string
	httpClient       *http.Client
	FileSaverService *FileSaverService
	Records          *Records
	apiKey           string
	allowGroups      []string
}

func NewMessageProcessorService(FileSaverService *FileSaverService, Line_webhook_url string, httpClient *http.Client, apiKey string, allowGroups []string) *MessageProcessorService {

	Records := &Records{
		latestSavedTimestamp: 0,
	}

	Records.LoadTimestampFromJsonfile()

	return &MessageProcessorService{
		Line_webhook_url: Line_webhook_url,
		httpClient:       httpClient,
		FileSaverService: FileSaverService,
		Records:          Records,
		apiKey:           apiKey,
		allowGroups:      allowGroups,
	}
}

func (mc *MessageProcessorService) SaveTimestampToJsonfile(timestamp int64) {
	os.WriteFile("latest_timestamp.json", []byte(fmt.Sprintf("%d", timestamp)), 0644)
}

func (r *Records) LoadTimestampFromJsonfile() error {
	data, err := os.ReadFile("latest_timestamp.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("latest_timestamp.json not found, initializing with 0")
			r.latestSavedTimestamp = 0
			return nil
		}
		return fmt.Errorf("failed to read latest_timestamp.json: %w", err)
	}
	var timestamp int64
	_, err = fmt.Sscanf(string(data), "%d", &timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp from latest_timestamp.json: %w", err)
	}
	r.latestSavedTimestamp = timestamp
	fmt.Println("Loaded latest saved timestamp:", r.latestSavedTimestamp)
	return nil
}

func (mc *MessageProcessorService) GetJsonArchiveLists() ([]string, error) {
	url := mc.Line_webhook_url + "/line_chat_webhook/list_filenames?more_than_timestamp=" + fmt.Sprintf("%d", mc.Records.latestSavedTimestamp)
	fmt.Println(url)

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+mc.apiKey) // Add the API key here

	// Send the request
	resp, err := mc.httpClient.Do(req)
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

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewReader(byteString))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mc.apiKey) // Add the API key here

	// Send the request
	resp, err := mc.httpClient.Do(req)
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
	req.Header.Set("Authorization", "Bearer "+mc.apiKey) // Add the API key here

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

func (s *Records) SaveLatestTimestamp(arr []string) {

	if len(arr) == 0 {
		fmt.Println("No timestamps to process")
		return
	}

	newArray := make([]int64, len(arr))
	for i, value := range arr {
		var timestamp int64
		_, err := fmt.Sscanf(value, "%d", &timestamp)
		if err != nil {
			fmt.Printf("Error parsing timestamp %s: %v\n", value, err)
			continue
		}
		newArray[i] = timestamp
	}

	maxValue := newArray[0]
	for _, value := range newArray {
		if value > maxValue {
			maxValue = value
		}
	}
	s.latestSavedTimestamp = maxValue
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (mc *MessageProcessorService) Process(MaximumFiles int) error {

	var selectListsFiles []string

	listsFiles, err := mc.GetJsonArchiveLists()
	if err != nil {

		return fmt.Errorf("failed to get JSON archive lists: %w", err)
	}

	if len(listsFiles) > MaximumFiles {
		selectListsFiles = listsFiles[:MaximumFiles]
	} else {
		selectListsFiles = listsFiles
	}

	MapFilenameChatdatas, err := mc.GetJsonArchives(selectListsFiles)
	if err != nil {
		return fmt.Errorf("failed to get JSON archives: %w", err)
	}

	for filename, chatDatas := range MapFilenameChatdatas {
		fmt.Println("FileName:", filename)
		for _, chatData := range chatDatas {
			for _, event := range chatData.Events {
				if event.Source.GroupID != "" && !contains(mc.allowGroups, event.Source.GroupID) {
					fmt.Println("Skipping event from group:", event.Source.GroupID)
					continue
				}
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
		// mc.DeleteFiles([]string{filename})
	}

	mc.Records.SaveLatestTimestamp(selectListsFiles)
	mc.SaveTimestampToJsonfile(mc.Records.latestSavedTimestamp)

	return nil
}
