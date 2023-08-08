package infra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"wfs/backend/apperr"
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

func (d *Discord) Upload(filename string, text string, message string) error {
	// create file
	file, err := os.Create(filename)
	if err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	defer func() {
		_ = file.Close()
		_ = os.Remove(filename)
	}()

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

	res, err := postMultipartFormData[any](d.config.URL, forms)
	if err != nil {
		return apperr.New(apperr.ErrDiscordAPI, err)
	}

	if res.StatusCode != http.StatusOK {
		//nolint:goerr113
		return apperr.New(apperr.ErrDiscordAPI, fmt.Errorf("status_code=%d, body=%s", res.StatusCode, res.BodyByte))
	}

	return nil
}
