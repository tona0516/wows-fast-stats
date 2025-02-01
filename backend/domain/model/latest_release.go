package model

type LatestRelease struct {
	Semver    string `json:"semver"`
	URL       string `json:"url"`
	Updatable bool   `json:"updatable"`
}
