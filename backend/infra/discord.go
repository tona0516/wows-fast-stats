package infra

import (
	"archive/zip"
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
	zipName := "out.zip"
	textName := "out.txt"

	// create zip
	file, err := os.Create(zipName)
	defer func() {
		_ = file.Close()
		_ = os.Remove(zipName)
	}()
	if err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	zw := zip.NewWriter(file)
	fw, err := zw.Create(textName)
	if err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	_, err = fw.Write([]byte(text))
	if err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	if err := zw.Close(); err != nil {
		return apperr.New(apperr.ErrWriteFile, err)
	}

	// upload zip
	//nolint:errchkjson
	payload, _ := json.Marshal(DisCordRequestBody{Content: message})
	forms := []Form{
		{name: "payload_json", content: string(payload), isFile: false},
		{name: "file", content: zipName, isFile: true},
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
