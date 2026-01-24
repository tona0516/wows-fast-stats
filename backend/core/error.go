package core

import "github.com/morikuni/failure"

const (
	ErrWGAPITemporaryUnavailable failure.StringCode = "wg_api_temporary_unavailable_error"
	ErrWGAPI                     failure.StringCode = "wg_api_error"
	ErrNumbersAPI                failure.StringCode = "numbers_api_error"
	ErrGithubAPI                 failure.StringCode = "github_api_error"
	ErrDiscordAPI                failure.StringCode = "discord_api_error"
	ErrClanAPI                   failure.StringCode = "clan_api_error"
	ErrWailsOpenDirectoryDialog  failure.StringCode = "wails_open_directory_dialog_error"
	ErrTempArenaInfoNotFound     failure.StringCode = "temp_arena_info_not_found_error"
	ErrTempArenaInfoSearch       failure.StringCode = "temp_arena_info_search_error"
	ErrJSONRead                  failure.StringCode = "json_read_error"
	ErrJSONNotFound              failure.StringCode = "json_not_found_error"
	ErrJSONWrite                 failure.StringCode = "json_write_error"
	ErrStringRead                failure.StringCode = "string_read_error"
	ErrStringWrite               failure.StringCode = "string_write_error"
	ErrInvalidExpectedStats      failure.StringCode = "invalid_expected_stats_error"
	ErrInvalidInstallPath        failure.StringCode = "invalid_install_path_error"
	ErrInitialSettingRequired    failure.StringCode = "initial_setting_required_error"
)
