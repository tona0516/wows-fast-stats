import type { core } from "@wails/go/models";
import type { ColorCode } from "./ColorCode";
import type { Tonako } from "./Tonako";

type BasicKey = "playerInfo" | "warship";
type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof core.PlayerStats, CommonMethod>;
export type ColumnCategory = Readonly<"basic" | StatsCategory>;

export type ShipType = Readonly<keyof core.ShipTypeGroup>;
export type StatsExtra = Exclude<keyof core.Player, BasicKey | CommonMethod>;

export type ShipStatsKey = Exclude<keyof core.ShipStats, CommonMethod>;
export type OverallStatsKey = Exclude<keyof core.OverallStats, CommonMethod>;
export type StatsKey = ShipStatsKey | OverallStatsKey;

export type Optional<T> = T | undefined;

export type Page = "stats" | "pref" | "info";

export type ColumnSettingPattern = "ship" | "overall" | "both";

export type PlayerNameColorType =
  | "shipPR"
  | "overallPR"
  | "threatLevel"
  | "none";

export type WarshipNamesColorType =
  | "shipType"
  | "none";

export type RowPattern =
  | "no_column"
  | "private"
  | "no_stats"
  | "no_ship_stats"
  | "full";

export type ColumnInfo = {
  readonly minName: string;
  readonly fullName: string;
};

export type StackedBarChartParam = {
  readonly label: string;
  readonly colorCode?: ColorCode;
  readonly value: number;
};

export type TonakoParam = {
  message: string;
  isLoading: boolean;
  tonako: Tonako;
};
