package vo

type LogParam struct {
	Timestamp string `json:"timestamp"`
	LogLevel  string `json:"log_level"`
	Semver    string `json:"semver"`
	Message   string `json:"message"`
	Error     string `json:"error"`
}
