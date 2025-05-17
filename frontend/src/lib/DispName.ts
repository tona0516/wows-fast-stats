import { ArrayMap } from "src/lib/ArrayMap";
import { FontSize } from "src/lib/FontSize";
import type { RatingLevel } from "src/lib/RatingLevel";
import type { SummaryShipType } from "src/lib/Summary";
import type {
  ColumnCategory,
  ShipType,
  StatsExtra,
  TierGroup,
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
    ["rank_solo", "ランク戦"],
  ]);

  export const COLUMN_CATEGORIES = new ArrayMap<ColumnCategory, string>([
    ["basic", "基本情報"],
    ["ship", "艦成績"],
    ["overall", "総合成績"],
  ]);

  export const SKILL_LEVELS = new ArrayMap<RatingLevel, string>([
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

  export const SHIP_TYPE_FOR_SUMMARY = new ArrayMap<SummaryShipType, string>([
    ...SHIP_TYPES,
    ["all", "全艦種"],
  ]);

  export const TIER_GROUPS = new ArrayMap<TierGroup, string>([
    ["low", "1~4"],
    ["middle", "5~7"],
    ["high", "8~★"],
  ]);

  export const MIN_COLUMN_NAMES = new ArrayMap<string, string>([
    ["pr", "PR"],
    ["damage", "Dmg"],
    ["max_damage", "最大Dmg"],
    ["win_rate", "勝率"],
    ["kd_rate", "K/D"],
    ["kill", "撃沈"],
    ["planes_killed", "撃墜"],
    ["exp", "Exp"],
    ["battles", "戦闘数"],
    ["survived_rate", "生存率(勝|敗)"],
    ["hit_rate", "命中率(主|魚)"],
    ["platoon_rate", "プラ比"],
    ["avg_tier", "平均T"],
    ["using_ship_type_rate", "艦種割合"],
    ["using_tier_rate", "T割合"],
    ["threat_level", "戦力評価"],
  ]);

  export const FULL_COLUMN_NAMES = new ArrayMap<string, string>([
    ["pr", "Personal Rating"],
    ["damage", "与ダメージ"],
    ["max_damage", "最大与ダメージ"],
    ["win_rate", "勝率"],
    ["kd_rate", "キルデス比"],
    ["kill", "撃沈"],
    ["planes_killed", "撃墜"],
    ["exp", "経験値(プレミアム補正含む)"],
    ["battles", "戦闘数"],
    ["survived_rate", "生存率(勝利|敗北)"],
    ["hit_rate", "命中率(主砲|魚雷)"],
    ["platoon_rate", "分艦隊比率"],
    ["avg_tier", "平均Tier"],
    ["using_ship_type_rate", "使用艦種割合"],
    ["using_tier_rate", "プレイTier割合"],
    ["threat_level", "戦力評価(闇深XVM算出ロジック)"],
  ]);
}
