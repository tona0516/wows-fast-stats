import type { vo } from "../wailsjs/go/models";
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

  export const DIGITS = {
    pr: 0,
    damage: 0,
    win_rate: 1,
    kd_rate: 2,
    kill: 2,
    survived_rate: 1,
    hit_rate: 1,
    planes_killed: 1,
    exp: 0,
    battles: 0,
    avg_tier: 2,
    using_ship_type_rate: 1,
    using_tier_rate: 1,
  };

  export const TYPE_S_COLORS = {
    cv: "#5E2883",
    bb: "#CA1028",
    cl: "#27853F",
    dd: "#D9760F",
    ss: "#233B8B",
  };

  export const TYPE_P_COLORS = {
    cv: "#CAB2D6",
    bb: "#FBB4C4",
    cl: "#CCEBC5",
    dd: "#FEE6AA",
    ss: "#B3CDE3",
  };

  export const TIER_S_COLORS = {
    low: "#8CA113",
    middle: "#205B85",
    high: "#990F4F",
  };

  export const TIER_P_COLORS = {
    low: "#E6F5B0",
    middle: "#B3D7DD",
    high: "#E3ADD5",
  };

  export const RANK_BG_COLORS = {
    bad: "#a41200",
    belowAvg: "#a34a02",
    avg: "#a38204",
    good: "#518517",
    veryGood: "#2f6f41",
    great: "#04436d",
    unicum: "#232166",
    superUnicum: "#531460",
  };

  export const RANK_TEXT_COLORS = {
    bad: "#ff382d",
    belowAvg: "#fd9234",
    avg: "#ffd351",
    good: "#57e500",
    veryGood: "#44b200",
    great: "#02f7da",
    unicum: "#da6ff5",
    superUnicum: "#bf15ee",
  };

  export const MAX_MEMO_LENGTH = 100;

  export const DEFAULT_ALERT_PLAYER: vo.AlertPlayer = {
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
}
