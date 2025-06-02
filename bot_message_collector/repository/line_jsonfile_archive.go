package repository

import (
	"bot_message_collector/entities"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type LineJsonfileArchive struct {
	listDataArchivePath string
}

func NewLineJsonfileArchive(listDataArchivePath string) *LineJsonfileArchive {
	return &LineJsonfileArchive{
		listDataArchivePath: listDataArchivePath,
	}
}

func (r *LineJsonfileArchive) SaveDataToJsonFile(ListData []entities.LineWebhook) error {

	byteData, err := json.Marshal(ListData)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%d.json", r.listDataArchivePath, time.Now().Unix())

	err = os.WriteFile(filename, byteData, 0644) // Clear the file before writing new data
	if err != nil {
		return err
	}

	return nil
}

func (r *LineJsonfileArchive) GetListFilenames(moreThanTimestampInt int) ([]string, error) {
	files, err := os.ReadDir(r.listDataArchivePath)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			if moreThanTimestampInt > 0 {
				fileTimestamp := file.Name()[:len(file.Name())-5]  // Remove the .json extension
				fileTimestampInt, _ := strconv.Atoi(fileTimestamp) // Remove the .json extension
				if fileTimestampInt > moreThanTimestampInt {
					filenames = append(filenames, file.Name())
				}
			} else {
				filenames = append(filenames, file.Name())
			}

		}
	}

	return filenames, nil
}

func (r LineJsonfileArchive) GetJsonArchives(filenames []string) (map[string][]entities.LineWebhook, error) {
	archives := make(map[string][]entities.LineWebhook)

	for _, filename := range filenames {
		r, err := r.loadDataFromJsonFile(filename)
		if err != nil {
			fmt.Printf("Error loading data from file %s: %v\n", filename, err)
			return nil, fmt.Errorf("error loading data from file %s: %w", filename, err)
		}
		archives[filename] = r
	}

	return archives, nil
}

func (r *LineJsonfileArchive) DeleteListTimestampLowerThan(timestamp int64) (deletedList []string, err error) {
	files, err := os.ReadDir(r.listDataArchivePath)
	if err != nil {
		return nil, err
	}

	deletedList = []string{}

	for _, file := range files {
		if !file.IsDir() {
			fileTimestampStr := file.Name()[:len(file.Name())-5] // Remove the .json extension
			fileTimestamp, err := strconv.ParseInt(fileTimestampStr, 10, 64)
			if err != nil {
				return deletedList, fmt.Errorf("error parsing timestamp from file name %s: %w", file.Name(), err)
			}

			if fileTimestamp < timestamp {
				err = os.Remove(fmt.Sprintf("%s/%s", r.listDataArchivePath, file.Name()))
				deletedList = append(deletedList, file.Name())
				if err != nil {
					return deletedList, fmt.Errorf("error deleting file %s: %w", file.Name(), err)
				}
			}
		}
	}

	return deletedList, nil
}

func (r *LineJsonfileArchive) loadDataFromJsonFile(filename string) ([]entities.LineWebhook, error) {

	filePath := fmt.Sprintf("%s/%s", r.listDataArchivePath, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var chatData []entities.LineWebhook
	err = json.Unmarshal(data, &chatData)
	if err != nil {
		return nil, err
	}

	return chatData, nil
}

func (r *LineJsonfileArchive) DeleteJsonFile(filenames []string) error {
	for _, filename := range filenames {
		filePath := fmt.Sprintf("%s/%s", r.listDataArchivePath, filename)

		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *LineJsonfileArchive) DeleteAllJsonFiles() error {
	files, err := os.ReadDir(r.listDataArchivePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			err := os.Remove(fmt.Sprintf("%s/%s", r.listDataArchivePath, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
