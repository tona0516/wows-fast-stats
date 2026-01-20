import { ColorCode } from "./ColorCode";
import type {
  ColumnCategory,
  ColumnInfo,
  PlayerNameColorPattern,
  ShipType,
  StatsKey,
} from "./types";

export const ZOOM_RATES = [
  50, 67, 75, 80, 90, 100, 110, 120, 125, 133, 150, 167, 175, 200,
] as const;

export const STATS_EXTRAS: Readonly<Map<string, string>> = new Map<
  string,
  string
>([
  ["pvpAll", "ランダム戦"],
  ["pvpSolo", "ランダム戦(ソロ)"],
  ["rankSolo", "ランク戦"],
]);

export const PLAYER_NAME_COLORS: Readonly<Map<PlayerNameColorPattern, string>> =
  new Map<PlayerNameColorPattern, string>([
    ["shipPR", "艦成績のPR"],
    ["overallPR", "総合成績のPR"],
    ["threatLevel", "戦力評価"],
    ["none", "なし"],
  ]);

export const STATS_COLUMN_INFO: {
  [key in StatsKey]: ColumnInfo;
} = {
  pr: { min: "PR", full: "Personal Rating", pattern: "both" },
  winRate: { min: "勝率", full: "勝率", pattern: "both" },
  damage: { min: "Dmg", full: "与ダメージ", pattern: "both" },
  maxDamage: { min: "最大Dmg", full: "最大与ダメージ", pattern: "both" },
  kdRate: { min: "K/D", full: "キルデス比", pattern: "both" },
  kill: { min: "撃沈", full: "撃沈", pattern: "both" },
  exp: { min: "Exp", full: "経験値(プレミアム補正含む)", pattern: "both" },
  battles: { min: "戦闘数", full: "戦闘数", pattern: "both" },
  platoonRate: { min: "分艦隊比", full: "分艦隊比率", pattern: "both" },
  avgTier: { min: "平均T", full: "平均Tier", pattern: "both" },
  planesKilled: { min: "撃墜", full: "撃墜", pattern: "ship" },
  survivedRate: {
    min: "生存率(全|勝|敗)",
    full: "生存率(全戦|勝利|敗北)",
    pattern: "ship",
  },
  hitRate: { min: "命中率(主|魚)", full: "命中率(主砲|魚雷)", pattern: "ship" },
  threatLevel: {
    min: "戦力評価",
    full: "戦力評価(闇深XVM算出ロジック)",
    pattern: "overall",
  },
  usingShipTypeRate: {
    min: "艦種割合",
    full: "使用艦種割合",
    pattern: "overall",
  },
  usingTierRate: {
    min: "T割合",
    full: "プレイTier割合",
    pattern: "overall",
  },
  efficiencyBadge: {
    min: "技能バッジ(E|1|2|3)",
    full: "技能バッジ(Expert|1st|2nd|3rd)",
    pattern: "both",
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

export const RATING_COLORS: { [key: string]: ColorCode } = {
  bad: new ColorCode("#FE0E00"),
  below_avg: new ColorCode("#FE7903"),
  avg: new ColorCode("#FFC71F"),
  good: new ColorCode("#44B300"),
  very_good: new ColorCode("#318000"),
  great: new ColorCode("#02C9B3"),
  unicum: new ColorCode("#D042F3"),
  super_unicum: new ColorCode("#A00DC5"),
} as const;

export const RATING_NAMES: { [key: string]: string } = {
  bad: "Bad",
  below_avg: "Below Average",
  avg: "Average",
  good: "Good",
  very_good: "Very Good",
  great: "Great",
  unicum: "Unicum",
  super_unicum: "Super Unicum",
} as const;

type ColorPair = {
  readonly text: ColorCode;
  readonly background: ColorCode;
};

export const THREAT_LEVEL_COLORS: { [rank: string]: ColorPair } = {
  ir: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#000000") },
  r: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#FF0000") },
  o: { text: new ColorCode("#331100"), background: new ColorCode("#FFA500") },
  y: { text: new ColorCode("#331100"), background: new ColorCode("#FFFF00") },
  g: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#008000") },
  b: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#2255FF") },
  i: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#234794") },
  v: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#705DA8") },
  uv: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#800080") },
} as const;
