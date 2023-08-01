import type { domain } from "../wailsjs/go/models";
import { Func, Page } from "./enums";

export namespace Const {
  export const BASE_NUMBERS_URL = "https://asia.wows-numbers.com/";

  export const COLUMN_NAMES = {
    // categoty
    basic: { min: "基本情報", full: "基本情報" },
    ship: { min: "艦成績", full: "艦成績" },
    overall: { min: "総合成績", full: "総合成績" },
    // value
    is_in_avg: { min: "", full: "" },
    player_name: { min: "プレイヤー", full: "プレイヤー" },
    ship_info: { min: "艦", full: "艦情報" },
    pr: { min: "PR", full: "Personal Rating" },
    damage: { min: "Dmg", full: "平均ダメージ" },
    win_rate: { min: "勝率", full: "勝率" },
    kd_rate: { min: "K/D", full: "K/D比" },
    kill: { min: "撃沈", full: "平均撃沈数" },
    planes_killed: { min: "撃墜", full: "平均撃墜数" },
    survived_rate: { min: "生存率(勝|負)", full: "生存率 (勝利|敗北)" },
    exp: { min: "Exp", full: "平均取得経験値" },
    battles: { min: "戦闘数", full: "戦闘数" },
    avg_tier: { min: "平均T", full: "平均Tier" },
    using_ship_type_rate: { min: "艦割合", full: "艦種別プレイ割合" },
    using_tier_rate: { min: "T割合", full: "ティア別プレイ割合" },
    hit_rate: { min: "Hit率(主|魚)", full: "命中率 (主砲|魚雷)" },
  };

  export const SKILL_LEVEL_LABELS = {
    bad: "Bad",
    below_avg: "Below Average",
    avg: "Average",
    good: "Good",
    very_good: "Very Good",
    great: "Great",
    unicum: "Unicum",
    super_unicum: "Super Unicum",
  };

  export const SHIP_TYPE_LABELS = {
    ss: "潜水艦",
    dd: "駆逐艦",
    cl: "巡洋艦",
    bb: "戦艦",
    cv: "空母",
  };

  export const TIER_GROUP_LABELS = {
    low: "1~4",
    middle: "5~7",
    high: "8~★",
  };

  export const MAX_MEMO_LENGTH = 100;

  export const DEFAULT_ALERT_PLAYER: domain.AlertPlayer = {
    account_id: 0,
    name: "",
    pattern: "bi-check-circle-fill",
    message: "",
  };

  export const PAGES: { title: string; name: Page; iconClass: string }[] = [
    {
      title: "ホーム",
      name: Page.Main,
      iconClass: "bi bi-house",
    },
    {
      title: "設定",
      name: Page.Config,
      iconClass: "bi bi-gear",
    },
    {
      title: "アプリ情報",
      name: Page.AppInfo,
      iconClass: "bi bi-info-circle",
    },
    {
      title: "プレイヤーリスト",
      name: Page.AlertPlayer,
      iconClass: "bi bi-person-lines-fill",
    },
  ];

  export const FUNCS: { title: string; name: Func; iconClass: string }[] = [
    {
      title: "リロード",
      name: Func.Reload,
      iconClass: "bi bi-arrow-clockwise",
    },
    {
      title: "スクリーンショット",
      name: Func.Screenshot,
      iconClass: "bi bi-camera",
    },
  ];

  export const FONT_SIZE = {
    "x-small": "極小",
    small: "小",
    medium: "中",
    large: "大",
    "x-large": "極大",
  };

  export const STATS_PATTERN = {
    pvp_solo: "ランダム戦(ソロ)",
    pvp_all: "ランダム戦",
  };

  export const PLAYER_NAME_COLOR = {
    ship: "艦成績のPR",
    overall: "総合成績のPR",
    none: "なし",
  };
}
