import type { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { data } from "wailsjs/go/models";

export const BasicKey = {
  player_info: "player_info",
  ship_info: "ship_info",
};
export type BasicKey = (typeof BasicKey)[keyof typeof BasicKey];

type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof data.PlayerStats, CommonMethod>;
export type ColumnCategory = Readonly<"basic" | StatsCategory>;

export type ShipType = Readonly<keyof data.ShipTypeGroup>;
export type TierGroup = Readonly<keyof data.TierGroup>;
export type StatsExtra = Exclude<
  keyof data.Player,
  keyof typeof BasicKey | CommonMethod
>;

export type ShipStatsKey = Exclude<keyof data.ShipStats, CommonMethod>;
export type OverallStatsKey = Exclude<keyof data.OverallStats, CommonMethod>;
export type StatsKey = ShipStatsKey | OverallStatsKey;

export type OptionalBattle = data.Battle | undefined;
export type OptionalTeamThreatLevels = TeamThreatLevel[] | undefined;

export type Page = "stats" | "ap_config" | "config" | "info";

export type GetStatsFunction = (ps: data.PlayerStats) => number;

export type ColumnSetting = { ship: boolean; overall: boolean; digit: number };

export type ColumnSettingPattern = "ship" | "overall" | "both";
