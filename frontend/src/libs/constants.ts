import { ColorCode } from "./ColorCode";
import type { ShipType, StatsKey } from "./types";

type ColumnInfo = {
  readonly minName: string;
  readonly fullName: string;
};

export const STATS_COLUMN_INFO: {
  [key in StatsKey]: ColumnInfo;
} = {
  pr: {
    minName: "PR",
    fullName: "Personal Rating",
  },
  winRate: {
    minName: "勝率",
    fullName: "勝率",
  },
  damage: {
    minName: "Dmg",
    fullName: "与ダメージ",
  },
  maxDamage: {
    minName: "最大Dmg",
    fullName: "最大与ダメージ",
  },
  kdRate: {
    minName: "K/D",
    fullName: "キルデス比",
  },
  kill: {
    minName: "撃沈",
    fullName: "撃沈",
  },
  exp: {
    minName: "Exp",
    fullName: "経験値(プレミアム補正含む)",
  },
  battles: {
    minName: "戦闘数",
    fullName: "戦闘数",
  },
  platoonRate: {
    minName: "分艦隊比",
    fullName: "分艦隊比率",
  },
  efficiencyBadge: {
    minName: "技能バッジ(E|1|2|3)",
    fullName: "技能バッジ(Expert|1st|2nd|3rd)",
  },
  planesKilled: {
    minName: "撃墜",
    fullName: "撃墜",
  },
  survivedRate: {
    minName: "生存率(全|勝|敗)",
    fullName: "生存率(全戦|勝利|敗北)",
  },
  hitRate: {
    minName: "命中率(主|魚)",
    fullName: "命中率(主砲|魚雷)",
  },
  avgTier: {
    minName: "平均T",
    fullName: "平均Tier",
  },
  threatLevel: {
    minName: "戦力評価",
    fullName: "戦力評価(闇深XVM算出ロジック)",
  },
  usingShipTypeRate: {
    minName: "艦種割合",
    fullName: "使用艦種割合",
  },
  usingTierRate: {
    minName: "T割合",
    fullName: "プレイTier割合",
  },
} as const;

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
