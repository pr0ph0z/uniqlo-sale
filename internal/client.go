package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Media struct {
	Type    string `json:"type"`
	Media   string `json:"media"`
	File    string `json:"file"`
	Caption string `json:"caption"`
}

type Request struct {
	ChatID string  `json:"chat_id"`
	Media  []Media `json:"media"`
}

func SendMediaGroup(medias []Media) (err error) {
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMediaGroup", botToken)

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	multipartWriter.WriteField("chat_id", chatID)

	mediaJSON, err := json.Marshal(medias)
	if err != nil {
		return
	}

	multipartWriter.WriteField("media", string(mediaJSON))

	for i, media := range medias {
		var (
			file *os.File
			part io.Writer
		)
		file, err = os.Open(media.File)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err = multipartWriter.CreateFormFile(fmt.Sprintf("file%d", i), media.File)
		if err != nil {
			return err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", baseURL, &requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		panic("failed to process the request")
	}

	return
}
