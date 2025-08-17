import { ArrayMap } from "src/lib/ArrayMap";
import type { RatingLevel } from "src/lib/RatingLevel";
import type { ShipType, TierGroup } from "src/lib/types";

export namespace DispName {
  export const SKILL_LEVELS = new Map<RatingLevel, string>([
    ["bad", "Bad"],
    ["below_avg", "Below Average"],
    ["avg", "Average"],
    ["good", "Good"],
    ["very_good", "Very Good"],
    ["great", "Great"],
    ["unicum", "Unicum"],
    ["super_unicum", "Super Unicum"],
  ]);

  export const SHIP_TYPES = new ArrayMap<ShipType, string>([
    ["ss", "潜水艦"],
    ["dd", "駆逐艦"],
    ["cl", "巡洋艦"],
    ["bb", "戦艦"],
    ["cv", "空母"],
  ]);

  export const TIER_GROUPS = new ArrayMap<TierGroup, string>([
    ["low", "1~4"],
    ["middle", "5~7"],
    ["high", "8~★"],
  ]);
}
