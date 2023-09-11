import { KeyValue } from "src/lib/KeyValue";
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
  export const STATS_PATTERNS: KeyValue<StatsExtra, string>[] = [
    new KeyValue("pvp_solo", "ランダム戦(ソロ)"),
    new KeyValue("pvp_all", "ランダム戦"),
  ];

  export const COLUMN_CATEGORIES: KeyValue<ColumnCategory, string>[] = [
    new KeyValue("basic", "基本情報"),
    new KeyValue("ship", "艦成績"),
    new KeyValue("overall", "総合成績"),
  ];

  export const SKILL_LEVELS: KeyValue<Rating, string>[] = [
    new KeyValue("bad", "Bad"),
    new KeyValue("below_avg", "Below Average"),
    new KeyValue("avg", "Average"),
    new KeyValue("good", "Good"),
    new KeyValue("very_good", "Very Good"),
    new KeyValue("great", "Great"),
    new KeyValue("unicum", "Unicum"),
    new KeyValue("super_unicum", "Super Unicum"),
  ];

  export const FONT_SIZES: KeyValue<FontSize, string>[] = [
    new KeyValue("x-small", "極小"),
    new KeyValue("small", "小"),
    new KeyValue("medium", "中"),
    new KeyValue("large", "大"),
    new KeyValue("x-large", "極大"),
  ];

  export const PLAYER_NAME_COLORS: KeyValue<PlayerNameColor, string>[] = [
    new KeyValue("ship", "艦成績のPR"),
    new KeyValue("overall", "総合成績のPR"),
    new KeyValue("none", "なし"),
  ];

  export const SHIP_TYPES: KeyValue<ShipType, string>[] = [
    new KeyValue("ss", "潜水艦"),
    new KeyValue("dd", "駆逐艦"),
    new KeyValue("cl", "巡洋艦"),
    new KeyValue("bb", "戦艦"),
    new KeyValue("cv", "空母"),
  ];

  export const TIER_GROUPS: KeyValue<TierGroup, string>[] = [
    new KeyValue("low", "1~4"),
    new KeyValue("middle", "5~7"),
    new KeyValue("high", "8~★"),
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
