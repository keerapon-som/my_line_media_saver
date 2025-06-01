package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/h2non/filetype"
)

type FileSaver struct {
	MessageID string `json:"message_id"`
	FileName  string `json:"file_name"`
}

type FileSaverService struct {
	accessToken      string
	mediaArchivePath string
	saveContentChan  chan FileSaver
}

func NewLineContentSaverService(accessToken string, mediaArchivePath string) *FileSaverService {
	return &FileSaverService{
		accessToken:      accessToken,
		mediaArchivePath: mediaArchivePath,
		saveContentChan:  make(chan FileSaver), // Buffered channel to handle multiple requests
	}
}

func (f *FileSaverService) SendToSaveContent(messageID string, filename string) {

	f.saveContentChan <- FileSaver{
		MessageID: messageID,
		FileName:  filename,
	}
}

func (f *FileSaverService) SaveContentWorker() {
	fmt.Println("current Len:", len(f.saveContentChan))

	for {
		select {
		case content := <-f.saveContentChan:
			fmt.Println("=== Start saving content for id  ===", content.MessageID)
			// === Request setup ===
			url := fmt.Sprintf("https://api-data.line.me/v2/bot/message/%s/content", content.MessageID)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", "Bearer "+f.accessToken)

			// === Send request ===
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			// === Check status ===
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("❌ Failed to download content: %d\n", resp.StatusCode)
				os.Exit(1)
			}

			// === Get content-type ===
			contentType := resp.Header.Get("Content-Type")
			fmt.Println("Content-Type:", contentType)

			// === MIME to extension map ===
			extMap := map[string]string{
				// "image/jpeg": ".jpg", "image/png": ".png", "image/gif": ".gif",
				// "video/mp4": ".mp4", "video/mpeg": ".mpeg", "video/quicktime": ".mov",
				// "audio/m4a": ".m4a", "audio/mpeg": ".mp3", "audio/wav": ".wav",
				// "application/pdf": ".pdf", "application/json": ".json",
				// "text/plain": ".txt", "text/csv": ".csv", "application/zip": ".zip",
			}

			ext, found := extMap[contentType]
			if !found {
				ext = ".bin" // fallback extension
			}

			// === Save file ===
			todayFolder := f.mediaArchivePath + "/" + time.Now().Format("02_01_2006")
			if _, err := os.Stat(todayFolder); os.IsNotExist(err) {
				err := os.MkdirAll(todayFolder, 0755)
				if err != nil {
					fmt.Println("Error creating folder:", err)
					return
				}
			}
			tempFilename := todayFolder + "/" + content.FileName + ext
			outFile, err := os.Create(tempFilename)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(outFile, resp.Body)
			if err != nil {
				panic(err)
			}
			outFile.Close() // Ensure file is closed after writing

			// === Fallback detection using filetype ===
			if ext == ".bin" {
				data, err := ioutil.ReadFile(tempFilename)
				if err != nil {
					fmt.Println("Error reading file:", err)
					return
				}

				kind, err := filetype.Match(data)
				if err != nil {
					fmt.Println("Error detecting file type:", err)
					return
				}

				if kind != filetype.Unknown {
					newExt := "." + kind.Extension
					newFilename := todayFolder + "/" + content.FileName + newExt
					err := os.Rename(tempFilename, newFilename)
					if err != nil {
						fmt.Println("Error renaming file:", err)
						return
					}
					fmt.Printf("✅ Detected type: %s, extension: %s\n", kind.MIME.Value, newExt)
				} else {
					fmt.Println("Unknown file type, saved as .bin")
				}
			} else {
				fmt.Printf("✅ Saved as: %s\n", content.FileName+ext)
			}

			// case <-time.After(5 * time.Second):
			// 	fmt.Println("hello")
		}
	}
}
