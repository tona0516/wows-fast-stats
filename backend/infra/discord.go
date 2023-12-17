package infra

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/infra/webapi"

	"github.com/morikuni/failure"
)

type Discord struct {
	config RequestConfig
}

type DiscordRequestBody struct {
	Content string `json:"content"`
}

func NewDiscord(config RequestConfig) *Discord {
	return &Discord{config: config}
}

func (d *Discord) Upload(filename string, text string, message string) error {
	// create file
	file, err := os.Create(filename)
	if err != nil {
		return failure.Wrap(err)
	}

	defer func() {
		_ = file.Close()
		_ = os.Remove(filename)
	}()

	// write text to file
	_, err = file.WriteString(text)
	if err != nil {
		return failure.Wrap(err)
	}

	// upload file
	//nolint:errchkjson
	payload, _ := json.Marshal(DiscordRequestBody{Content: message})
	forms := []webapi.Form{
		webapi.NewForm("payload_json", string(payload), false),
		webapi.NewForm("file", filename, true),
	}

	res, err := webapi.PostMultipartFormData[any](d.config.URL, d.config.Timeout, forms)
	if err != nil {
		return failure.Wrap(err)
	}

	if res.StatusCode != http.StatusOK {
		return failure.New(
			apperr.DiscordAPISendLogError,
			failure.Context{
				"status_code": strconv.Itoa(res.StatusCode),
				"body":        string(res.ByteBody),
			},
		)
	}

	return nil
}
