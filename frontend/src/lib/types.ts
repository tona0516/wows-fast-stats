import type { Summary } from "src/lib/Summary";
import { model } from "wailsjs/go/models";

export const BasicKey = {
  player_info: "player_info",
  ship_info: "ship_info",
};
export type BasicKey = (typeof BasicKey)[keyof typeof BasicKey];

type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof model.PlayerStats, CommonMethod>;
export type ColumnCategory = Readonly<"basic" | StatsCategory>;

export type ShipType = Readonly<keyof model.ShipTypeGroup>;
export type TierGroup = Readonly<keyof model.TierGroup>;
export type StatsExtra = Exclude<
  keyof model.Player,
  keyof typeof BasicKey | CommonMethod
>;
export type Rating = Readonly<keyof model.SkillColorCode>;

export type ShipKey = Readonly<keyof model.Ship>;
export type OverallKey = Readonly<keyof model.Overall>;
export type CommonKey = ShipKey & OverallKey;
export type DigitKey = Readonly<keyof model.CustomDigit>;

export type OptionalBattle = model.Battle | undefined;
export type OptionalSummary = Summary | undefined;
