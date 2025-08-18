import type { data } from "wailsjs/go/models";

type BasicKey = "player_info" | "ship_info";
type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof data.PlayerStats, CommonMethod>;
export type ColumnCategory = Readonly<"basic" | StatsCategory>;

export type ShipType = Readonly<keyof data.ShipTypeGroup>;
export type TierGroup = Readonly<keyof data.TierGroup>;
export type StatsExtra = Exclude<keyof data.Player, BasicKey | CommonMethod>;

export type ShipStatsKey = Exclude<keyof data.ShipStats, CommonMethod>;
export type OverallStatsKey = Exclude<keyof data.OverallStats, CommonMethod>;
export type StatsKey = ShipStatsKey | OverallStatsKey;

export type OptionalBattle = data.Battle | undefined;

export type Page = "stats" | "ap_config" | "config" | "info";

export type ColumnSettingPattern = "ship" | "overall" | "both";

export type PlayerNameColor = "ship" | "overall" | "none";

export type ThreatLevel = "ir" | "r" | "o" | "y" | "g" | "b" | "i" | "v" | "uv";

export type Rating =
  | "bad"
  | "below_avg"
  | "avg"
  | "good"
  | "very_good"
  | "great"
  | "unicum"
  | "super_unicum";

export type RowPattern =
  | "no_column"
  | "private"
  | "no_stats"
  | "no_ship_stats"
  | "full";

export type ColumnSetting = {
  readonly ship: boolean;
  readonly overall: boolean;
  readonly digit: number;
};

export type ThreatLevelInfo = {
  readonly level: ThreatLevel;
  readonly threshold: number;
};

export type TeamThreatLevel = {
  readonly average: number;
  readonly dissociationDegree: number;
  readonly accuracy: number;
};

export type RatingInfo = {
  readonly rating: Rating;
  readonly prThreshold: number;
  readonly damageThreshold: number;
  readonly winRateThreshold: number;
};

export type ColumnInfo = {
  readonly min: string;
  readonly full: string;
  readonly pattern: ColumnSettingPattern;
};

export type ColorPair = {
  readonly text: string;
  readonly background: string;
};
