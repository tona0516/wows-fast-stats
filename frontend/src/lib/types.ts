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
export type ShipKey = Readonly<keyof data.UCDisplayShip>;
export type OverallKey = Readonly<keyof data.UCDisplayOverall>;
export type CommonKey = ShipKey & OverallKey;
export type DigitKey = Readonly<keyof data.UCDigit>;

export type OptionalBattle = data.Battle | undefined;
export type OptionalTeamThreatLevels = TeamThreatLevel[] | undefined;

export type Page = "stats" | "ap_config" | "config" | "info";

export type GetStatsFunction = (ps: data.PlayerStats) => number;

export type PlayerDetail = {
  id: number;
  name: string;
  clan?: {
    id: number;
    tag: string;
  };
};
