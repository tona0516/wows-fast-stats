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
  TD_STR: "td-string",
  OMIT: "omit",
};
export type CssClass = (typeof CssClass)[keyof typeof CssClass];

export const BasicKey = {
  player_info: "player_info",
  ship_info: "ship_info",
};
export type BasicKey = (typeof BasicKey)[keyof typeof BasicKey];

export const ScreenshotType = {
  manual: "manual",
  auto: "auto",
};
export type ScreenshotType =
  (typeof ScreenshotType)[keyof typeof ScreenshotType];

type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof domain.PlayerStats, CommonMethod>;
export type ColumnCategory = Readonly<"basic" | StatsCategory>;
export type ShipType = Readonly<keyof domain.ShipTypeGroup>;
export type TierGroup = Readonly<keyof domain.TierGroup>;
export type StatsExtra = Exclude<
  keyof domain.Player,
  keyof typeof BasicKey | CommonMethod
>;
export type Rating = Readonly<keyof domain.SkillColorCode>;

export type ShipKey = Readonly<keyof domain.Ship>;
const ships = Object.keys(new domain.Ship());
export const includesShips = (key: string): boolean => {
  return ships.includes(key);
};

export type OverallKey = Readonly<keyof domain.Overall>;
const overalls = Object.keys(new domain.Overall());
export const includesOveralls = (key: string): boolean => {
  return overalls.includes(key);
};

export type CommonKey = ShipKey & OverallKey;

export type DigitKey = Readonly<keyof domain.CustomDigit>;

export type OptionalBattle = domain.Battle | undefined;
export type OptionalSummary = Summary | undefined;
