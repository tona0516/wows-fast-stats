package vo

type ShipTypeGroup[T any] struct {
    SS T `json:"ss"`
    DD T `json:"dd"`
    CL T `json:"cl"`
    BB T `json:"bb"`
    CV T `json:"cv"`
}
