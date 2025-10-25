package domain

type NewVersion struct {
	Semver string `json:"semver"`
	URL    string `json:"url"`
}
