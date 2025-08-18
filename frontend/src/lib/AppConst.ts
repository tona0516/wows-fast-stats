import type {
  ColorPair,
  ColumnCategory,
  ColumnInfo,
  PlayerNameColor,
  Rating,
  RatingInfo,
  ShipType,
  StatsExtra,
  StatsKey,
  ThreatLevel,
  ThreatLevelInfo,
  TierGroup,
} from "./types";

export namespace AppConst {
  export const ZOOM_RATES = [
    25, 33, 50, 67, 75, 80, 90, 100, 110, 120, 125, 133, 150, 167, 175, 200,
    250, 300, 400, 500,
  ] as const;

  export const STATS_EXTRAS: Readonly<Map<StatsExtra, string>> = new Map<
    StatsExtra,
    string
  >([
    ["pvp_all", "ランダム戦"],
    ["pvp_solo", "ランダム戦(ソロ)"],
    ["rank_solo", "ランク戦"],
  ]);

  export const LIGHTER_THEMES: string[] = ["light", "nord", "garden"] as const;
  export const DARKER_THEMES: string[] = ["dark", "night", "black"] as const;
  export const THEMES: string[] = [
    ...LIGHTER_THEMES,
    ...DARKER_THEMES,
  ] as const;

  export const PLAYER_NAME_COLORS: Readonly<Map<PlayerNameColor, string>> =
    new Map<PlayerNameColor, string>([
      ["ship", "艦成績のPR"],
      ["overall", "総合成績のPR"],
      ["none", "なし"],
    ]);

  export const STATS_COLUMN_INFO: {
    [key in StatsKey]: ColumnInfo;
  } = {
    pr: { min: "PR", full: "Personal Rating", pattern: "both" },
    win_rate: { min: "勝率", full: "勝率", pattern: "both" },
    damage: { min: "Dmg", full: "与ダメージ", pattern: "both" },
    max_damage: { min: "最大Dmg", full: "最大与ダメージ", pattern: "both" },
    kd_rate: { min: "K/D", full: "キルデス比", pattern: "both" },
    kill: { min: "撃沈", full: "撃沈", pattern: "both" },
    exp: { min: "Exp", full: "経験値(プレミアム補正含む)", pattern: "both" },
    battles: { min: "戦闘数", full: "戦闘数", pattern: "both" },
    platoon_rate: { min: "分艦隊比", full: "分艦隊比率", pattern: "both" },
    avg_tier: { min: "平均T", full: "平均Tier", pattern: "both" },
    planes_killed: { min: "撃墜", full: "撃墜", pattern: "ship" },
    survived_rate: {
      min: "生存率",
      full: "生存率(勝利|敗北)",
      pattern: "ship",
    },
    hit_rate: { min: "命中率", full: "命中率(主砲|魚雷)", pattern: "ship" },
    threat_level: {
      min: "戦力評価",
      full: "戦力評価(闇深XVM算出ロジック)",
      pattern: "overall",
    },
    using_ship_type_rate: {
      min: "艦種割合",
      full: "使用艦種割合",
      pattern: "overall",
    },
    using_tier_rate: {
      min: "T割合",
      full: "プレイTier割合",
      pattern: "overall",
    },
  } as const;

  export const STATS_KEYS = Object.keys(
    STATS_COLUMN_INFO,
  ) as readonly StatsKey[];

  export const CATEGORY_NAMES: Readonly<Map<ColumnCategory, string>> = new Map<
    ColumnCategory,
    string
  >([
    ["basic", "基本情報"],
    ["ship", "艦成績"],
    ["overall", "総合成績"],
  ]);

  export const SKILL_LEVELS: Readonly<Map<Rating, string>> = new Map<
    Rating,
    string
  >([
    ["bad", "Bad"],
    ["below_avg", "Below Average"],
    ["avg", "Average"],
    ["good", "Good"],
    ["very_good", "Very Good"],
    ["great", "Great"],
    ["unicum", "Unicum"],
    ["super_unicum", "Super Unicum"],
  ]);

  export const SHIP_TYPES: Readonly<Map<ShipType, string>> = new Map<
    ShipType,
    string
  >([
    ["ss", "潜水艦"],
    ["dd", "駆逐艦"],
    ["cl", "巡洋艦"],
    ["bb", "戦艦"],
    ["cv", "空母"],
  ]);

  export const TIER_GROUPS: Readonly<Map<TierGroup, string>> = new Map<
    TierGroup,
    string
  >([
    ["low", "1~4"],
    ["middle", "5~7"],
    ["high", "8~★"],
  ]);

  export const NUMBERS_URL = "https://asia.wows-numbers.com/";

  export const SHIP_TYPE_COLORS: Readonly<Map<ShipType, string>> = new Map<
    ShipType,
    string
  >([
    ["ss", "#233B8B"],
    ["dd", "#D9760F"],
    ["cl", "#27853F"],
    ["bb", "#CA1028"],
    ["cv", "#5E2883"],
  ]);

  export const TIER_GROUP_COLORS: Readonly<Map<TierGroup, string>> = new Map<
    TierGroup,
    string
  >([
    ["low", "#8CA113"],
    ["middle", "#205B85"],
    ["high", "#990F4F"],
  ]);

  export const THREAT_LEVEL_COLORS: Readonly<Map<ThreatLevel, ColorPair>> =
    new Map<ThreatLevel, ColorPair>([
      ["ir", { text: "#FFFFFF", background: "#000000" }],
      ["r", { text: "#FFFFFF", background: "#FF0000" }],
      ["o", { text: "#331100", background: "#FFA500" }],
      ["y", { text: "#331100", background: "#FFFF00" }],
      ["g", { text: "#FFFFFF", background: "#008000" }],
      ["b", { text: "#FFFFFF", background: "#2255FF" }],
      ["i", { text: "#FFFFFF", background: "#234794" }],
      ["v", { text: "#FFFFFF", background: "#705DA8" }],
      ["uv", { text: "#FFFFFF", background: "#800080" }],
    ]);

  export const RATING_COLORS: Readonly<Map<Rating, string>> = new Map<
    Rating,
    string
  >([
    ["bad", "#FE0E00"],
    ["below_avg", "#FE7903"],
    ["avg", "#FFC71F"],
    ["good", "#44B300"],
    ["very_good", "#318000"],
    ["great", "#02C9B3"],
    ["unicum", "#D042F3"],
    ["super_unicum", "#A00DC5"],
  ]);

  const THREAT_LEVEL_COEF = 0.5;
  export const THREAT_LEVELS_INFO: ThreatLevelInfo[] = [
    { level: "ir", threshold: 0 },
    { level: "r", threshold: 8000 * THREAT_LEVEL_COEF },
    { level: "o", threshold: 13000 * THREAT_LEVEL_COEF },
    { level: "y", threshold: 19000 * THREAT_LEVEL_COEF },
    { level: "g", threshold: 25000 * THREAT_LEVEL_COEF },
    { level: "b", threshold: 32000 * THREAT_LEVEL_COEF },
    { level: "i", threshold: 35000 * THREAT_LEVEL_COEF },
    { level: "v", threshold: 40000 * THREAT_LEVEL_COEF },
    { level: "uv", threshold: 44000 * THREAT_LEVEL_COEF },
  ] as const;

  export const RATING_INFO: RatingInfo[] = [
    { rating: "bad", prThreshold: 0, damageThreshold: 0, winRateThreshold: 0 },
    {
      rating: "below_avg",
      prThreshold: 750,
      damageThreshold: 0.6,
      winRateThreshold: 47,
    },
    {
      rating: "avg",
      prThreshold: 1100,
      damageThreshold: 0.8,
      winRateThreshold: 50,
    },
    {
      rating: "good",
      prThreshold: 1350,
      damageThreshold: 1.0,
      winRateThreshold: 52,
    },
    {
      rating: "very_good",
      prThreshold: 1550,
      damageThreshold: 1.2,
      winRateThreshold: 54,
    },
    {
      rating: "great",
      prThreshold: 1750,
      damageThreshold: 1.4,
      winRateThreshold: 56,
    },
    {
      rating: "unicum",
      prThreshold: 2100,
      damageThreshold: 1.5,
      winRateThreshold: 60,
    },
    {
      rating: "super_unicum",
      prThreshold: 2450,
      damageThreshold: 1.6,
      winRateThreshold: 65,
    },
  ] as const;

  export const ROMAN_NUMERALS: { [key: number]: string } = {
    1: "I",
    2: "II",
    3: "III",
    4: "IV",
    5: "V",
    6: "VI",
    7: "VII",
    8: "VIII",
    9: "IX",
    10: "X",
  } as const;
}
