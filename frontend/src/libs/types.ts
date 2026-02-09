import type { core } from "@wails/go/models";
import type { ColorCode } from "./ColorCode";
import type { Tonako } from "./Tonako";

type BasicKey = "playerInfo" | "warship";
type CommonMethod = "convertValues";

export type StatsCategory = Exclude<keyof core.PlayerStats, CommonMethod>;
export type ShipType = Readonly<keyof core.ShipTypeGroup>;
export type StatsExtra = Exclude<keyof core.Player, BasicKey | CommonMethod>;

type ShipStatsKey = Exclude<keyof core.ShipStats, CommonMethod>;
type OverallStatsKey = Exclude<keyof core.OverallStats, CommonMethod>;
export type StatsKey = ShipStatsKey | OverallStatsKey;

export type Optional<T> = T | undefined;

export type Page = "stats" | "setting" | "info";

export type PlayerNameColorType = "shipPR" | "overallPR" | "none";

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
