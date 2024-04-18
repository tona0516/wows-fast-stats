package data

type GHLatestRelease struct {
	TagName   string `json:"tag_name"`
	HTMLURL   string `json:"html_url"`
	Updatable bool   `json:"updatable"`
}
