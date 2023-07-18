package vo

type ValidatedResult struct {
	InstallPath string `json:"install_path"`
	AppID       string `json:"appid"`
}

func (r *ValidatedResult) Valid() bool {
	return r.InstallPath == "" && r.AppID == ""
}
