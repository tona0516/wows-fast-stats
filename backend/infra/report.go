package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wfs/backend/application/vo"
)

const emptyJSON = "{}"

type Report struct {
	env       vo.Env
	localFile LocalFile
	discord   Discord
	storage   Storage
}

func NewReport(
	env vo.Env,
	localFile LocalFile,
	discord Discord,
	storage Storage,
) *Report {
	return &Report{
		env:       env,
		localFile: localFile,
		discord:   discord,
		storage:   storage,
	}
}

func (r *Report) Send(message string, content error) {
	// get UserConfig
	userConfig, err := r.storage.ReadUserConfig()
	if err != nil || !userConfig.SendReport {
		return
	}
	jsonUserConfig := prettryJSON(userConfig)

	// get IGN
	ign, _ := r.storage.ReadOwnIGN()

	// get TempArenaInfo
	tempArenaInfo, err := r.localFile.TempArenaInfo(userConfig.InstallPath)
	var jsonTempArenaInfo string
	if err != nil {
		jsonTempArenaInfo = emptyJSON
	} else {
		jsonTempArenaInfo = prettryJSON(tempArenaInfo)
	}

	// make content
	targets := []string{
		"Semver:",
		fmt.Sprintf("%+v\n", r.env.Semver),
		"IGN:",
		fmt.Sprintf("%+v\n", ign),
		"Error:",
		fmt.Sprintf("%+v\n", content),
		"TempArenaInfo:",
		jsonTempArenaInfo,
		"UserConfig:",
		jsonUserConfig,
	}

	sb := strings.Builder{}
	for _, target := range targets {
		sb.WriteString(target)
		sb.WriteString("\n")
	}

	// send report
	_ = r.discord.Upload("report.txt", sb.String(), message)
}

func prettryJSON(data any) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return emptyJSON
	}

	return buffer.String()
}
