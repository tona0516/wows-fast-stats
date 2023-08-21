package vo

// Note: set by ldflags in Makefile.
type Env struct {
	AppName string `json:"app_name"`
	IsDev   bool   `json:"is_dev"`
	Semver  string `json:"semver"`
}
