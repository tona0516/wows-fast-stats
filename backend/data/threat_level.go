package data

type ThreatLevel struct {
	Rank     ThreatLevelRank `json:"rank"`
	Raw      float64         `json:"raw"`
	Modified float64         `json:"modified"`
}

type ThreatLevelRank string

const (
	ThreatLevelRankUV ThreatLevelRank = "uv"
	ThreatLevelRankV  ThreatLevelRank = "v"
	ThreatLevelRankI  ThreatLevelRank = "i"
	ThreatLevelRankB  ThreatLevelRank = "b"
	ThreatLevelRankG  ThreatLevelRank = "g"
	ThreatLevelRankY  ThreatLevelRank = "y"
	ThreatLevelRankO  ThreatLevelRank = "o"
	ThreatLevelRankR  ThreatLevelRank = "r"
	ThreatLevelRankIR ThreatLevelRank = "ir"
)
