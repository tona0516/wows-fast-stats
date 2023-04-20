package vo

type Displays struct {
    PR bool `json:"pr"`
    ShipDamage bool `json:"ship_damage"`
    ShipWinRate bool `json:"ship_win_rate"`
    ShipKdRate bool `json:"ship_kd_rate"`
    ShipBattles bool `json:"ship_battles"`
    PlayerDamage bool `json:"player_damage"`
    PlayerWinRate bool `json:"player_win_rate"`
    PlayerKdRate bool `json:"player_kd_rate"`
    PlayerBattles bool `json:"player_battles"`
    PlayerAvgTier bool `json:"player_avg_tier"`
}

type UserConfig struct {
    InstallPath string `json:"install_path"`
    Appid string `json:"appid"`
    FontSize string `json:"font_size"`
    Displays Displays `json:"displays"`
}
