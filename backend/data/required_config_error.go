package data

type RequiredConfigError struct {
	Valid       bool   `json:"valid"`
	InstallPath string `json:"install_path"`
	AppID       string `json:"appid"`
}
