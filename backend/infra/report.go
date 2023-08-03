package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wfs/backend/application/vo"
)

type Report struct {
	env       vo.Env
	localFile LocalFile
	discord   Discord
}

func NewReport(
	env vo.Env,
	localFile LocalFile,
	discord Discord,
) *Report {
	return &Report{
		env:       env,
		localFile: localFile,
		discord:   discord,
	}
}

func (r *Report) Send(content any) {
	// get UserConfig
	userConfig, err := r.localFile.User()
	if err != nil || !userConfig.SendReport {
		return
	}
	jsonUserConfig := prettryJSON(userConfig)

	// get TempArenaInfo
	tempArenaInfo, err := r.localFile.TempArenaInfo(userConfig.InstallPath)
	jsonTempArenaInfo := "{}"
	if err == nil {
		jsonTempArenaInfo = prettryJSON(tempArenaInfo)
	}

	// make content
	targets := []string{
		"Semver:",
		fmt.Sprintf("%+v\n", r.env.Semver),
		"Content:",
		fmt.Sprintf("%+v\n", content),
		"UserConfig:",
		jsonUserConfig,
		"TempArenaInfo:",
		jsonTempArenaInfo,
	}

	sb := strings.Builder{}
	for _, target := range targets {
		sb.WriteString(target)
		sb.WriteString("\n")
	}

	// send report
	message := "uploaded file!"
	if r.env.IsDebug {
		message += " [dev]"
	}

	_ = r.discord.Upload(sb.String(), message)
}

func prettryJSON(data any) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return "{}"
	}

	return buffer.String()
}
