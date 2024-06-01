package data

type WGClansInfo map[int]WGClansInfoData

func (w WGClansInfo) ToArray() []WGClansInfoData {
	array := make([]WGClansInfoData, 0)
	for _, v := range w {
		array = append(array, v)
	}
	return array
}

type WGClansInfoData struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}
