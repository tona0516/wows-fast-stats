import type { data } from "@wails/go/models";

type BasicKey = "player_info" | "ship_info";
type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof data.PlayerStats, CommonMethod>;
export type ColumnCategory = Readonly<"basic" | StatsCategory>;

export type ShipType = Readonly<keyof data.ShipTypeGroup>;
export type StatsExtra = Exclude<keyof data.Player, BasicKey | CommonMethod>;

export type ShipStatsKey = Exclude<keyof data.ShipStats, CommonMethod>;
export type OverallStatsKey = Exclude<keyof data.OverallStats, CommonMethod>;
export type StatsKey = ShipStatsKey | OverallStatsKey;

export type Optional<T> = T | undefined;

export type Page = "stats" | "config" | "info";

export type ColumnSettingPattern = "ship" | "overall" | "both";

export type PlayerNameColorPattern =
  | "pr_ship"
  | "pr_overall"
  | "threat_level"
  | "none";

export type RowPattern =
  | "no_column"
  | "private"
  | "no_stats"
  | "no_ship_stats"
  | "full";

export type ColumnInfo = {
  readonly min: string;
  readonly full: string;
  readonly pattern: ColumnSettingPattern;
};

export type TierGroup = Readonly<keyof data.TierGroup>;

export type Rating =
  | "bad"
  | "below_avg"
  | "avg"
  | "good"
  | "very_good"
  | "great"
  | "unicum"
  | "super_unicum";
