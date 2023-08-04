package infra

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"wfs/backend/apperr"

	"github.com/pkg/errors"
)

type Discord struct {
	config RequestConfig
}

type DisCordRequestBody struct {
	Content string `json:"content"`
}

func NewDiscord(config RequestConfig) *Discord {
	return &Discord{config: config}
}

func (d *Discord) Upload(text string, message string) error {
	filename := "out.txt"

	// create file
	file, err := os.Create(filename)
	defer func() {
		_ = file.Close()
		_ = os.Remove(filename)
	}()
	if err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	// write text to file
	_, err = file.WriteString(text)
	if err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	// upload file
	//nolint:errchkjson
	payload, _ := json.Marshal(DisCordRequestBody{Content: message})
	forms := []Form{
		{name: "payload_json", content: string(payload), isFile: false},
		{name: "file", content: filename, isFile: true},
	}

	response, err := postMultipartFormData[any](d.config.URL, forms, d.config.Retry)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		message := "status_code: " +
			strconv.Itoa(response.StatusCode) +
			" response_body: " +
			response.BodyString
		return apperr.New(apperr.ErrDiscordAPI, errors.New(message))
	}

	return nil
}
