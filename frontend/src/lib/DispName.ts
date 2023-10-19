import { ArrayMap } from "src/lib/ArrayMap";
import { FontSize } from "src/lib/FontSize";
import {
  type ColumnCategory,
  type Rating,
  type ShipType,
  type StatsExtra,
  type TierGroup,
} from "src/lib/types";

enum PlayerNameColor {
  SHIP = "ship",
  OVERALL = "overall",
  NONE = "none",
}

export namespace DispName {
  export const STATS_PATTERNS = new ArrayMap<StatsExtra, string>([
    ["pvp_all", "ランダム戦"],
    ["pvp_solo", "ランダム戦(ソロ)"],
  ]);

  export const COLUMN_CATEGORIES = new ArrayMap<ColumnCategory, string>([
    ["basic", "基本情報"],
    ["ship", "艦成績"],
    ["overall", "総合成績"],
  ]);

  export const SKILL_LEVELS = new ArrayMap<Rating, string>([
    ["bad", "Bad"],
    ["below_avg", "Below Average"],
    ["avg", "Average"],
    ["good", "Good"],
    ["very_good", "Very Good"],
    ["great", "Great"],
    ["unicum", "Unicum"],
    ["super_unicum", "Super Unicum"],
  ]);

  export const FONT_SIZES = new ArrayMap<FontSize, string>([
    [FontSize.XX_LARGE, "3"],
    [FontSize.X_LARGE, "2"],
    [FontSize.LARGE, "1"],
    [FontSize.MEDIUM, "0(デフォルト)"],
    [FontSize.SMALL, "-1"],
    [FontSize.X_SMALL, "-2"],
    [FontSize.XX_SMALL, "-3"],
  ]);

  export const PLAYER_NAME_COLORS = new ArrayMap<PlayerNameColor, string>([
    [PlayerNameColor.SHIP, "艦成績のPR"],
    [PlayerNameColor.OVERALL, "総合成績のPR"],
    [PlayerNameColor.NONE, "なし"],
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
