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
import { type KeyValue } from "src/lib/value_object/KeyValue";
import { type NavItem } from "src/lib/value_object/NavItem";

export namespace DispName {
  export const STATS_PATTERNS: KeyValue<StatsExtra, string>[] = [
    { key: "pvp_solo", value: "ランダム戦(ソロ)" },
    { key: "pvp_all", value: "ランダム戦" },
  ];

  export const COLUMN_CATEGORIES: KeyValue<ColumnCategory, string>[] = [
    { key: "basic", value: "基本情報" },
    { key: "ship", value: "艦成績" },
    { key: "overall", value: "総合成績" },
  ];

  export const SKILL_LEVELS: KeyValue<Rating, string>[] = [
    { key: "bad", value: "Bad" },
    { key: "below_avg", value: "Below Average" },
    { key: "avg", value: "Average" },
    { key: "good", value: "Good" },
    { key: "very_good", value: "Very Good" },
    { key: "great", value: "Great" },
    { key: "unicum", value: "Unicum" },
    { key: "super_unicum", value: "Super Unicum" },
  ];

  export const FONT_SIZES: KeyValue<FontSize, string>[] = [
    { key: "x-small", value: "極小" },
    { key: "small", value: "小" },
    { key: "medium", value: "中" },
    { key: "large", value: "大" },
    { key: "x-large", value: "極大" },
  ];

  export const PLAYER_NAME_COLORS: KeyValue<PlayerNameColor, string>[] = [
    { key: "ship", value: "艦成績のPR" },
    { key: "overall", value: "総合成績のPR" },
    { key: "none", value: "なし" },
  ];

  export const SHIP_TYPES: KeyValue<ShipType, string>[] = [
    { key: "ss", value: "潜水艦" },
    { key: "dd", value: "駆逐艦" },
    { key: "cl", value: "巡洋艦" },
    { key: "bb", value: "戦艦" },
    { key: "cv", value: "空母" },
  ];

  export const TIER_GROUPS: KeyValue<TierGroup, string>[] = [
    { key: "low", value: "1~4" },
    { key: "middle", value: "5~7" },
    { key: "high", value: "8~★" },
  ];

  export const PAGES: NavItem<Page>[] = [
    { type: Page.MAIN, dispName: "ホーム", icon: "bi bi-house" },
    { type: Page.CONFIG, dispName: "設定", icon: "bi bi-gear" },
    {
      type: Page.ALERT_PLAYER,
      dispName: "プレイヤーリスト",
      icon: "bi bi-person-lines-fill",
    },
    { type: Page.APPINFO, dispName: "アプリ情報", icon: "bi bi-info-circle" },
  ];

  export const FUNCS: NavItem<Func>[] = [
    { type: Func.RELOAD, dispName: "リロード", icon: "bi bi-arrow-clockwise" },
    {
      type: Func.SCREENSHOT,
      dispName: "スクリーンショット",
      icon: "bi bi-camera",
    },
  ];
}
