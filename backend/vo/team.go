package vo

type Team struct {
    Players Players `json:"players"`
    Name string `json:"name"`
}
