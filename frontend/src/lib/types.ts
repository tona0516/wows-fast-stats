import type { Summary } from "src/lib/Summary";
import { domain } from "wailsjs/go/models";

const FONT_SIZES = ["x-small", "small", "medium", "large", "x-large"] as const;
export type FontSize = (typeof FONT_SIZES)[number];

const PLAYER_NAME_COLORS = ["ship", "overall", "none"] as const;
export type PlayerNameColor = (typeof PLAYER_NAME_COLORS)[number];

export const Page = {
  MAIN: "main",
  CONFIG: "config",
  APPINFO: "appinfo",
  ALERT_PLAYER: "updatable",
};
export type Page = (typeof Page)[keyof typeof Page];

export const Func = {
  RELOAD: "reload",
  SCREENSHOT: "screenshot",
};
export type Func = (typeof Func)[keyof typeof Func];

export const ToastKey = {
  NEED_CONFIG: "need_config",
  WAIT: "wait",
  FETCHING: "fetching",
  UPDATABLE: "updatable",
  ERROR: "error",
};
export type ToastKey = (typeof ToastKey)[keyof typeof ToastKey];

// Note: see watcher.go
export const AppEvent = {
  BATTLE_START: "BATTLE_START",
  BATTLE_END: "BATTLE_END",
  BATTLE_ERR: "BATTLE_ERR",
  LOG: "LOG",
  ONLOAD: "ONLOAD",
};
export type AppEvent = (typeof AppEvent)[keyof typeof AppEvent];

export const RowPattern = {
  NO_COLUMN: "no_column",
  PRIVATE: "private",
  NO_DATA: "no_data",
  NO_SHIP_STATS: "no_ship_stats",
  FULL: "full",
};
export type RowPattern = (typeof RowPattern)[keyof typeof RowPattern];

export const CssClass = {
  TD_NUM: "td-number",
  TD_MULTI: "td-multiple",
};
export type CssClass = (typeof CssClass)[keyof typeof CssClass];

type CommonMethod = "convertValues";

export type ColumnCategory = Exclude<keyof domain.Displays, CommonMethod>;
export type StatsCategory = Exclude<keyof domain.PlayerStats, CommonMethod>;
export type ShipType = Readonly<keyof domain.ShipTypeGroup>;
export type TierGroup = Readonly<keyof domain.TierGroup>;
export type StatsExtra = Exclude<
  keyof domain.Player,
  "player_info" | "ship_info" | CommonMethod
>;
export type Rating = Readonly<keyof domain.SkillColorCode>;

export type BasicKey = Readonly<keyof domain.Basic>;
export type ShipKey = Readonly<keyof domain.Ship>;
export type OverallKey = Readonly<keyof domain.Overall>;

export type CommonStatsKey = ShipKey & OverallKey;
export type ShipOnlyKey = Exclude<ShipKey, CommonStatsKey>;
export type OverallOnlyKey = Exclude<OverallKey, CommonStatsKey>;

export type OptionalBattle = domain.Battle | undefined;
export type OptionalSummary = Summary | undefined;
