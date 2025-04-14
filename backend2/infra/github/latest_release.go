package github

type LatestRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}
