import { PlayerNameColor } from "./enums";
import type {
  ColumnCategory,
  ColumnSettingPattern,
  StatsExtra,
  StatsKey,
} from "./types";

export namespace AppConstants {
  export const ZOOM_RATES = [
    25, 33, 50, 67, 75, 80, 90, 100, 110, 120, 125, 133, 150, 167, 175, 200,
    250, 300, 400, 500,
  ];

  export const STATS_EXTRAS = new Map<StatsExtra, string>([
    ["pvp_all", "ランダム戦"],
    ["pvp_solo", "ランダム戦(ソロ)"],
    ["rank_solo", "ランク戦"],
  ]);

  export const LIGHTER_THEMES = ["light", "nord", "garden"];
  export const DARKER_THEMES = ["dark", "night", "black"];
  export const THEMES = [...LIGHTER_THEMES, ...DARKER_THEMES];

  export const PLAYER_NAME_COLORS = new Map<PlayerNameColor, string>([
    [PlayerNameColor.SHIP, "艦成績のPR"],
    [PlayerNameColor.OVERALL, "総合成績のPR"],
    [PlayerNameColor.NONE, "なし"],
  ]);

  class ColumnInfo {
    constructor(
      public readonly min: string,
      public readonly full: string,
      public readonly pattern: ColumnSettingPattern,
    ) {}
  }

  export const STATS_COLUMN_INFO: {
    [key in StatsKey]: ColumnInfo;
  } = {
    threat_level: new ColumnInfo(
      "戦力評価",
      "戦力評価(闇深XVM算出ロジック)",
      "overall",
    ),
    pr: new ColumnInfo("PR", "Personal Rating", "both"),
    win_rate: new ColumnInfo("勝率", "勝率", "both"),
    damage: new ColumnInfo("Dmg", "与ダメージ", "both"),
    max_damage: new ColumnInfo("最大Dmg", "最大与ダメージ", "both"),
    kd_rate: new ColumnInfo("K/D", "キルデス比", "both"),
    kill: new ColumnInfo("撃沈", "撃沈", "both"),
    exp: new ColumnInfo("Exp", "経験値(プレミアム補正含む)", "both"),
    battles: new ColumnInfo("戦闘数", "戦闘数", "both"),
    planes_killed: new ColumnInfo("撃墜", "撃墜", "ship"),
    survived_rate: new ColumnInfo("生存率", "生存率(勝利|敗北)", "ship"),
    hit_rate: new ColumnInfo("命中率", "命中率(主砲|魚雷)", "ship"),
    platoon_rate: new ColumnInfo("分艦隊比", "分艦隊比率", "overall"),
    avg_tier: new ColumnInfo("平均T", "平均Tier", "ship"),
    using_ship_type_rate: new ColumnInfo("艦種割合", "使用艦種割合", "overall"),
    using_tier_rate: new ColumnInfo("T割合", "プレイTier割合", "overall"),
  };

  export const STATS_KEYS = Object.keys(STATS_COLUMN_INFO) as StatsKey[];

  export const CATEGORY_NAMES = new Map<ColumnCategory, string>([
    ["basic", "基本情報"],
    ["ship", "艦成績"],
    ["overall", "総合成績"],
  ]);
}
