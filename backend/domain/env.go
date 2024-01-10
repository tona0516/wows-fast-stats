package domain

// Note: set by ldflags.
type Env struct {
	AppName string `json:"app_name"`
	IsDev   bool   `json:"is_dev"`
	Semver  string `json:"semver"`
}
