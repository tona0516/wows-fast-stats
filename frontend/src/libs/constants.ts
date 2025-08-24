import type { data } from "@wails/go/models";
import { ColorCode } from "./ColorCode";
import type {
  ColumnCategory,
  ColumnInfo,
  PlayerNameColor,
  ShipType,
  StatsExtra,
  StatsKey,
} from "./types";

export const ZOOM_RATES = [
  25, 33, 50, 67, 75, 80, 90, 100, 110, 120, 125, 133, 150, 167, 175, 200, 250,
  300, 400, 500,
] as const;

export const STATS_EXTRAS: Readonly<Map<StatsExtra, string>> = new Map<
  StatsExtra,
  string
>([
  ["pvp_all", "ランダム戦"],
  ["pvp_solo", "ランダム戦(ソロ)"],
  ["rank_solo", "ランク戦"],
]);

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

export const STATS_KEYS = Object.keys(STATS_COLUMN_INFO) as readonly StatsKey[];

export const CATEGORY_NAMES: Readonly<Map<ColumnCategory, string>> = new Map<
  ColumnCategory,
  string
>([
  ["basic", "基本情報"],
  ["ship", "艦成績"],
  ["overall", "総合成績"],
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

export const SHIP_TYPE_COLORS: Readonly<Map<ShipType, ColorCode>> = new Map<
  ShipType,
  ColorCode
>([
  ["ss", new ColorCode("#233B8B")],
  ["dd", new ColorCode("#D9760F")],
  ["cl", new ColorCode("#27853F")],
  ["bb", new ColorCode("#CA1028")],
  ["cv", new ColorCode("#5E2883")],
]);

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

export const DEFAULT_ALERT_PLAYER = {
  account_id: 0,
  name: "",
  pattern: "bi-check-circle-fill",
  message: "",
} as data.AlertPlayer;
