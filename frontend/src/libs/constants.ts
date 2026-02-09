import { ColorCode } from "./ColorCode";
import type { ShipType, StatsKey } from "./types";

type ColumnInfo = {
  readonly minName: string;
  readonly discription: string;
};

export const STATS_COLUMN_INFO: {
  [key in StatsKey]: ColumnInfo;
} = {
  pr: {
    minName: "PR",
    discription:
      "WoWS NumberのPersonal Rating\n詳細は以下を参照してください\nhttps://asia.wows-numbers.com/personal/rating",
  },
  winRate: {
    minName: "勝率",
    discription: "勝率(%)",
  },
  damage: {
    minName: "ダメージ",
    discription: "1戦あたりの平均与ダメージ",
  },
  maxDamage: {
    minName: "最大ダメージ",
    discription:
      "艦成績の場合: 使用艦種の最大与ダメージ\n総合成績の場合: すべての艦艇での最大与ダメージと艦名",
  },
  kdRate: {
    minName: "K/D",
    discription: "キルデス比",
  },
  kill: {
    minName: "撃沈",
    discription: "1戦あたりの平均撃沈数",
  },
  exp: {
    minName: "経験値",
    discription: "プレミアムアカウント補正された1戦あたりの平均取得経験値",
  },
  battles: {
    minName: "戦闘数",
    discription: "戦闘数",
  },
  platoonRate: {
    minName: "分艦隊比",
    discription:
      "戦闘数に対する分艦隊戦闘数の割合(1~3)\n1ならすべてソロ、3ならすべて3人分艦隊を意味します",
  },
  efficiencyBadge: {
    minName: "技能バッジ",
    discription:
      "艦成績の場合: 使用艦艇が保持する技能バッジ\n総合成績の場合: Expert、1st、2nd、3rdの各バッジ獲得数",
  },
  planesKilled: {
    minName: "撃墜",
    discription: "1戦あたりの平均撃墜数",
  },
  survivedRate: {
    minName: "生存率",
    discription: "戦闘数/勝利戦闘数/敗北戦闘数に対する生存率(%)",
  },
  hitRate: {
    minName: "命中率",
    discription: "主砲および魚雷の命中率(%)",
  },
  avgTier: {
    minName: "平均ティア",
    discription: "出撃した戦闘の平均ティア(1~11)",
  },
  usingShipTypeRate: {
    minName: "艦種割合",
    discription: "各艦種で出撃した戦闘数の割合(%)",
  },
  usingTierRate: {
    minName: "T割合",
    discription: "1~4/5~7/8~★のティアにおける戦闘数の割合(%)",
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
