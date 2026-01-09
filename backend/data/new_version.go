package data

type NewVersion struct {
	Version     string `json:"semver"`
	DownloadURL string `json:"url"`
}
