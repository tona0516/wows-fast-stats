import { ArrayMap } from "src/lib/ArrayMap";
import type { Summary } from "src/lib/Summary";
import { domain } from "wailsjs/go/models";

export const BASE_NUMBERS_URL = "https://asia.wows-numbers.com/";
export const MAX_MEMO_LENGTH = 100;
export const MAIN_PAGE_ID = "mainpage";

export const ADD_ALERT_PLAYER_MODAL_ID = "add-alert-player-modal";
export const EDIT_ALERT_PLAYER_MODAL_ID = "edit-alert-player-modal";
export const REMOVE_ALERT_PLAYER_MODAL_ID = "remove-alert-player-modal";
export const CONFIRM_MODAL_ID = "confirm-modal";

export const EMPTY_ALERT_PLAYER: domain.AlertPlayer = {
  account_id: 0,
  name: "",
  pattern: "bi-check-circle-fill",
  message: "",
} as const;

const FONT_SIZES = [
  "xx-small",
  "x-small",
  "small",
  "medium",
  "large",
  "x-large",
  "xx-large",
] as const;
export type FontSize = (typeof FONT_SIZES)[number];
export const ZOOM_RATIO = new ArrayMap<FontSize, number>([
  ["xx-small", 0.55],
  ["x-small", 0.7],
  ["small", 0.85],
  ["medium", 1.0],
  ["large", 1.15],
  ["x-large", 1.3],
  ["xx-large", 1.55],
]);
const PLAYER_NAME_COLORS = ["ship", "overall", "none"] as const;
export type PlayerNameColor = (typeof PLAYER_NAME_COLORS)[number];

export const Page = {
  MAIN: "main",
  CONFIG: "config",
  APPINFO: "appinfo",
};
export type Page = (typeof Page)[keyof typeof Page];

export const Func = {
  RELOAD: "reload",
  SCREENSHOT: "screenshot",
};
export type Func = (typeof Func)[keyof typeof Func];

export const ToastKey = {
  NEED_CONFIG: "need_config",
  UPDATABLE: "updatable",
  ERROR: "error",
};
export type ToastKey = (typeof ToastKey)[keyof typeof ToastKey];

// Note: see watcher.go
export const AppEvent = {
  BATTLE_START: "BATTLE_START",
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
