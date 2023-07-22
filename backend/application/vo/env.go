package vo

// Note: set by ldflags in Makefile.
type Env struct {
	AppName string `json:"app_name"`
	IsDebug bool   `json:"is_debug"`
	Semver  string `json:"semver"`
}
