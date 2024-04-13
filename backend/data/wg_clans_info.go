package data

type WGClansInfo map[int]WGClansInfoData

func (w WGClansInfo) Tags() []string {
	tags := make([]string, 0)
	for _, v := range w {
		tags = append(tags, v.Tag)
	}
	return tags
}

type WGClansInfoData struct {
	Tag string `json:"tag"`
}
