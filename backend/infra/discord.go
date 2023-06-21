package infra

import (
	"archive/zip"
	"changeme/backend/apperr"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type Discord struct {
	apiClient APIClientInterface[any]
}

type DisCordRequestBody struct {
	Content string `json:"content"`
}

func NewDiscord(webhookURL string) *Discord {
	return &Discord{
		apiClient: NewAPIClient[any](webhookURL),
	}
}

func (d *Discord) Upload(text string) (APIResponse[any], error) {
	response := APIResponse[any]{}
	zipName := "out.zip"
	textName := "out.txt"

	// create zip
	file, err := os.Create(zipName)
	defer func() {
		_ = file.Close()
		_ = os.Remove(zipName)
	}()
	if err != nil {
		return response, apperr.New(apperr.WriteFile, err)
	}

	zw := zip.NewWriter(file)
	fw, err := zw.Create(textName)
	if err != nil {
		return response, apperr.New(apperr.WriteFile, err)
	}

	_, err = fw.Write([]byte(text))
	if err != nil {
		return response, apperr.New(apperr.WriteFile, err)
	}

	if err := zw.Close(); err != nil {
		return response, apperr.New(apperr.WriteFile, err)
	}

	// upload zip
	//nolint:errchkjson
	payload, _ := json.Marshal(DisCordRequestBody{Content: "uploaded file!"})
	forms := []Form{
		{name: "payload_json", content: string(payload), isFile: false},
		{name: "file", content: zipName, isFile: true},
	}

	response, err = d.apiClient.PostMultipartFormData(forms)
	if err != nil {
		return response, err
	}

	if response.StatusCode != http.StatusOK {
		message := "status_code: " +
			strconv.Itoa(response.StatusCode) +
			" response_body: " +
			response.BodyString
		return response, apperr.New(apperr.DiscordAPIError, errors.New(message))
	}

	return response, nil
}
