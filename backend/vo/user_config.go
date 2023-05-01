package vo

type Basic struct {
    IsContainedAvg bool `json:"is_contained_avg"`
    PlaterName bool `json:"player_name"`
    ShipInfo bool `json:"ship_info"`
}

type Ship struct {
    PR bool `json:"pr"`
    Damage bool `json:"damage"`
    WinRate bool `json:"win_rate"`
    KdRate bool `json:"kd_rate"`
    WinSurvivedRate bool `json:"win_survived_rate"`
    LoseSurvivedRate bool `json:"lose_survived_rate"`
    Exp bool `json:"exp"`
    Battles bool `json:"battles"`
}

type Overall struct {
    Damage bool `json:"damage"`
    WinRate bool `json:"win_rate"`
    KdRate bool `json:"kd_rate"`
    WinSurvivedRate bool `json:"win_survived_rate"`
    LoseSurvivedRate bool `json:"lose_survived_rate"`
    Exp bool `json:"exp"`
    Battles bool `json:"battles"`
    AvgTier bool `json:"avg_tier"`
    UsingShipTypeRate bool `json:"using_ship_type_rate"`
    UsingTierRate bool `json:"using_tier_rate"`
}

type Displays struct {
    Basic Basic `json:"basic"`
    Ship Ship `json:"ship"`
    Overall Overall `json:"overall"`
}

type UserConfig struct {
    InstallPath string `json:"install_path"`
    Appid string `json:"appid"`
    FontSize string `json:"font_size"`
    Displays Displays `json:"displays"`
    SaveScreenshot bool `json:"save_screenshot"`
    SaveTempArenaInfo bool `json:"save_temp_arena_info"`
}
