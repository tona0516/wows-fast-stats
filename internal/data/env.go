package data

// Note: set by ldflags.
type Env struct {
	AppName string `json:"app_name"`
	WGAppID string `json:"app_id"`
	IsDev   bool   `json:"is_dev"`
	Semver  string `json:"semver"`
}
