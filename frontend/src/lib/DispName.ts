import { ArrayMap } from "src/lib/ArrayMap";
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
import { type NavItem } from "src/lib/value_object/NavItem";

export namespace DispName {
  export const STATS_PATTERNS = new ArrayMap<StatsExtra, string>([
    ["pvp_all", "ランダム戦"],
    ["pvp_solo", "ランダム戦(ソロ)"],
  ]);

  export const COLUMN_CATEGORIES = new ArrayMap<ColumnCategory, string>([
    ["basic", "基本情報"],
    ["ship", "艦成績"],
    ["overall", "総合成績"],
  ]);

  export const SKILL_LEVELS = new ArrayMap<Rating, string>([
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
    ["x-small", "極小"],
    ["small", "小"],
    ["medium", "中"],
    ["large", "大"],
    ["x-large", "極大"],
  ]);

  export const PLAYER_NAME_COLORS = new ArrayMap<PlayerNameColor, string>([
    ["ship", "艦成績のPR"],
    ["overall", "総合成績のPR"],
    ["none", "なし"],
  ]);

  export const SHIP_TYPES = new ArrayMap<ShipType, string>([
    ["ss", "潜水艦"],
    ["dd", "駆逐艦"],
    ["cl", "巡洋艦"],
    ["bb", "戦艦"],
    ["cv", "空母"],
  ]);

  export const TIER_GROUPS = new ArrayMap<TierGroup, string>([
    ["low", "1~4"],
    ["middle", "5~7"],
    ["high", "8~★"],
  ]);

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
