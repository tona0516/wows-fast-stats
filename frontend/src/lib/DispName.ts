import { Pair } from "src/lib/Pair";
import { Triple } from "src/lib/Triple";
import {
  Func,
  Page,
  type ColumnCategory,
  type FontSize,
  type PlayerNameColor,
  type Rating,
  type ShipType,
  type StatsExtra,
  type TierGroup,
} from "src/lib/types";

export namespace DispName {
  export const STATS_PATTERNS: Pair<StatsExtra, string>[] = [
    new Pair("pvp_solo", "ランダム戦(ソロ)"),
    new Pair("pvp_all", "ランダム戦"),
  ];

  export const COLUMN_CATEGORIES: Pair<ColumnCategory, string>[] = [
    new Pair("basic", "基本情報"),
    new Pair("ship", "艦成績"),
    new Pair("overall", "総合成績"),
  ];

  export const SKILL_LEVELS: Pair<Rating, string>[] = [
    new Pair("bad", "Bad"),
    new Pair("below_avg", "Below Average"),
    new Pair("avg", "Average"),
    new Pair("good", "Good"),
    new Pair("very_good", "Very Good"),
    new Pair("great", "Great"),
    new Pair("unicum", "Unicum"),
    new Pair("super_unicum", "Super Unicum"),
  ];

  export const FONT_SIZES: Pair<FontSize, string>[] = [
    new Pair("x-small", "極小"),
    new Pair("small", "小"),
    new Pair("medium", "中"),
    new Pair("large", "大"),
    new Pair("x-large", "極大"),
  ];

  export const PLAYER_NAME_COLORS: Pair<PlayerNameColor, string>[] = [
    new Pair("ship", "艦成績のPR"),
    new Pair("overall", "総合成績のPR"),
    new Pair("none", "なし"),
  ];

  export const SHIP_TYPES: Pair<ShipType, string>[] = [
    new Pair("ss", "潜水艦"),
    new Pair("dd", "駆逐艦"),
    new Pair("cl", "巡洋艦"),
    new Pair("bb", "戦艦"),
    new Pair("cv", "空母"),
  ];

  export const TIER_GROUPS: Pair<TierGroup, string>[] = [
    new Pair("low", "1~4"),
    new Pair("middle", "5~7"),
    new Pair("high", "8~★"),
  ];

  export const PAGES: Triple<Page, string, string>[] = [
    new Triple(Page.MAIN, "ホーム", "bi bi-house"),
    new Triple(Page.CONFIG, "設定", "bi bi-gear"),
    new Triple(Page.APPINFO, "アプリ情報", "bi bi-info-circle"),
    new Triple(
      Page.ALERT_PLAYER,
      "プレイヤーリスト",
      "bi bi-person-lines-fill",
    ),
  ];

  export const FUNCS: Triple<Func, string, string>[] = [
    new Triple(Func.RELOAD, "リロード", "bi bi-arrow-clockwise"),
    new Triple(Func.SCREENSHOT, "スクリーンショット", "bi bi-camera"),
  ];
}
