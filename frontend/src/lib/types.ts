import type { Summary } from "src/lib/Summary";
import { domain } from "wailsjs/go/models";

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
export type OverallKey = Readonly<keyof domain.Overall>;
export type CommonKey = ShipKey & OverallKey;
export type DigitKey = Readonly<keyof domain.CustomDigit>;

export type OptionalBattle = domain.Battle | undefined;
export type OptionalSummary = Summary | undefined;
