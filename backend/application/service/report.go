package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
)

type Report struct {
	version   vo.Version
	localFile repository.LocalFileInterface
	discord   repository.DiscordInterface
}

func NewReport(
	version vo.Version,
	localFile repository.LocalFileInterface,
	discord repository.DiscordInterface,
) *Report {
	return &Report{
		version:   version,
		localFile: localFile,
		discord:   discord,
	}
}

func (r *Report) Send(content error) error {
	// get UserConfig
	userConfig, err := r.localFile.User()
	if err != nil {
		return err
	}
	if !userConfig.SendReport {
		return nil
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
		fmt.Sprintf("%+v\n", r.version.Semver),
		"Revision:",
		fmt.Sprintf("%+v\n", r.version.Revision),
		"Error:",
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
	return r.discord.Upload(sb.String())
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
