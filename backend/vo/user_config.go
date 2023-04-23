package vo

type Displays struct {
    PlaterName bool `json:"player_name"`
    ShipInfo bool `json:"ship_info"`
    PR bool `json:"pr"`
    ShipDamage bool `json:"ship_damage"`
    ShipWinRate bool `json:"ship_win_rate"`
    ShipKdRate bool `json:"ship_kd_rate"`
    ShipWinSurvivedRate bool `json:"ship_win_survived_rate"`
    ShipLoseSurvivedRate bool `json:"ship_lose_survived_rate"`
    ShipExp bool `json:"ship_exp"`
    ShipBattles bool `json:"ship_battles"`
    PlayerDamage bool `json:"player_damage"`
    PlayerWinRate bool `json:"player_win_rate"`
    PlayerKdRate bool `json:"player_kd_rate"`
    PlayerWinSurvivedRate bool `json:"player_win_survived_rate"`
    PlayerLoseSurvivedRate bool `json:"player_lose_survived_rate"`
    PlayerExp bool `json:"player_exp"`
    PlayerBattles bool `json:"player_battles"`
    PlayerAvgTier bool `json:"player_avg_tier"`
    PlayerUsingShipTypeRate bool `json:"player_using_ship_type_rate"`
    PlayerUsingTierRate bool `json:"player_using_tier_rate"`
}

type UserConfig struct {
    InstallPath string `json:"install_path"`
    Appid string `json:"appid"`
    FontSize string `json:"font_size"`
    Displays Displays `json:"displays"`
    SaveScreenshot bool `json:"save_screenshot"`
}
