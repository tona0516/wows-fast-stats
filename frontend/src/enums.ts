export enum Page {
  Main = "main",
  Config = "config",
  AppInfo = "appinfo",
  AlertPlayer = "alert_player",
}

export enum Func {
  Reload = "reload",
  Screenshot = "screenshot",
}

/**
 * private: hidden player.
 * nodata: invalid player(bot/deleted account) or no battle for all random battle.
 * noshipstats: no battles in for random battle with the ship.
 * full: all values exists.
 */
export enum DisplayPattern {
  Private = "private",
  NoData = "nodata",
  NoShipStats = "noshipstats",
  Full = "full",
}

export enum StatsCategory {
  Basic = "basic",
  Ship = "ship",
  Overall = "overall",
}

export enum SkillLevel {
  Bad = "bad",
  BelowAvg = "belowAvg",
  Avg = "avg",
  Good = "good",
  VeryGood = "veryGood",
  Great = "great",
  Unicum = "unicum",
  SuperUnicum = "superUnicum",
}
