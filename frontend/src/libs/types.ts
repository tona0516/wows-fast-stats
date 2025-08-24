import type { data } from "@wails/go/models";
import type { ColorCode } from "./ColorCode";
import type { Tonako } from "./Tonako";

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

export type Page = "stats" | "ap_config" | "config" | "info";

export type ColumnSettingPattern = "ship" | "overall" | "both";

export type PlayerNameColor = "ship" | "overall" | "none";

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

export type TeamThreatLevel = {
  readonly average: number;
  readonly dissociationDegree: number;
  readonly accuracy: number;
};

export type ColumnInfo = {
  readonly min: string;
  readonly full: string;
  readonly pattern: ColumnSettingPattern;
};

export type StackedBarChartParam = {
  readonly label: string;
  readonly colorCode?: ColorCode;
  readonly value: number;
};

export type EditModalMode = "create" | "specify" | "edit";
export type EditModalParam = {
  mode: EditModalMode;
  form: data.AlertPlayer;
};

export type TonakoParam = {
  message: string;
  isLoading: boolean;
  tonako: Tonako;
};
