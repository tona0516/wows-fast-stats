package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wfs/backend/infra"
	"wfs/backend/vo"
)

type Report struct {
	discord           infra.DiscordInterface
	configRepo        infra.ConfigInterface
	tempArenaInfoRepo infra.TempArenaInfoInterface
}

func NewReport(
	discord infra.DiscordInterface,
	configRepo infra.ConfigInterface,
	tempArenaInfoRepo infra.TempArenaInfoInterface,
) *Report {
	return &Report{
		discord:           discord,
		configRepo:        configRepo,
		tempArenaInfoRepo: tempArenaInfoRepo,
	}
}

func (r *Report) Send(version vo.Version, content error) error {
	// get UserConfig
	userConfig, err := r.configRepo.User()
	if err != nil {
		return err
	}
	if !userConfig.SendReport {
		return nil
	}
	jsonUserConfig := prettryJSON(userConfig)

	// get TempArenaInfo
	tempArenaInfo, err := r.tempArenaInfoRepo.Get(userConfig.InstallPath)
	jsonTempArenaInfo := "{}"
	if err == nil {
		jsonTempArenaInfo = prettryJSON(tempArenaInfo)
	}

	// make content
	targets := []string{
		"Semver:",
		fmt.Sprintf("%+v\n", version.Semver),
		"Revision:",
		fmt.Sprintf("%+v\n", version.Revision),
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
	_, err = r.discord.Upload(sb.String())
	if err != nil {
		return err
	}

	return nil
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
