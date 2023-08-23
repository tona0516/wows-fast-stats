import iconCommonWealth from "src/assets/images/nation-commonwealth.png";
import iconEurope from "src/assets/images/nation-europe.png";
import iconFrance from "src/assets/images/nation-france.png";
import iconGermany from "src/assets/images/nation-germany.png";
import iconItaly from "src/assets/images/nation-italy.png";
import iconJapan from "src/assets/images/nation-japan.png";
import iconNetherlands from "src/assets/images/nation-netherlands.png";
import iconNone from "src/assets/images/nation-none.png";
import iconPanAmerica from "src/assets/images/nation-pan-america.png";
import iconPanAsia from "src/assets/images/nation-pan-asia.png";
import iconSpain from "src/assets/images/nation-spain.png";
import iconUk from "src/assets/images/nation-uk.png";
import iconUsa from "src/assets/images/nation-usa.png";
import iconUssr from "src/assets/images/nation-ussr.png";
import { Func, Page, SkillLevel } from "src/enums";
import type { domain } from "wailsjs/go/models";

export class ColumnName {
  min: string;
  full: string;

  constructor(min: string, full: string) {
    this.min = min;
    this.full = full;
  }
}

export class NavigationItem<T> {
  name: T;
  title: string;
  iconClass: string;

  constructor(name: T, title: string, iconClass: string) {
    this.title = title;
    this.name = name;
    this.iconClass = iconClass;
  }
}

export class SkillLevelItem {
  level: SkillLevel;
  tier: number;
  shipType: string;
  minPR: number;
  maxPR: number;
  minDamage: number;
  maxDamage: number;
  minWin: number;
  maxWin: number;

  constructor(
    level: SkillLevel,
    tier: number,
    shipType: string,
    minPR: number,
    maxPR: number,
    minDamage: number,
    maxDamage: number,
    minWin: number,
    maxWin: number
  ) {
    this.level = level;
    this.tier = tier;
    this.shipType = shipType;
    this.minPR = minPR;
    this.maxPR = maxPR;
    this.minDamage = minDamage;
    this.maxDamage = maxDamage;
    this.minWin = minWin;
    this.maxWin = maxWin;
  }
}

export namespace Const {
  export const BASE_NUMBERS_URL = "https://asia.wows-numbers.com/";

  export const COLUMN_NAMES = {
    // categoty
    basic: new ColumnName("基本情報", "基本情報"),
    ship: new ColumnName("艦成績", "艦成績"),
    overall: new ColumnName("総合成績", "総合成績"),
    // value
    is_in_avg: new ColumnName("", ""),
    player_name: new ColumnName("プレイヤー", "プレイヤー"),
    ship_info: new ColumnName("艦", "艦情報"),
    pr: new ColumnName("PR", "Personal Rating"),
    damage: new ColumnName("Dmg", "平均ダメージ"),
    win_rate: new ColumnName("勝率", "勝率"),
    kd_rate: new ColumnName("K/D", "K/D比"),
    kill: new ColumnName("撃沈", "平均撃沈数"),
    planes_killed: new ColumnName("撃墜", "平均撃墜数"),
    survived_rate: new ColumnName("生存率(勝|負)", "生存率 (勝利|敗北)"),
    exp: new ColumnName("Exp", "平均取得経験値"),
    battles: new ColumnName("戦闘数", "戦闘数"),
    avg_tier: new ColumnName("平均T", "平均Tier"),
    using_ship_type_rate: new ColumnName("艦割合", "艦種別プレイ割合"),
    using_tier_rate: new ColumnName("T割合", "ティア別プレイ割合"),
    hit_rate: new ColumnName("Hit率(主|魚)", "命中率 (主砲|魚雷)"),
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

  export const PAGES: NavigationItem<Page>[] = [
    new NavigationItem(Page.Main, "ホーム", "bi bi-house"),
    new NavigationItem(Page.Config, "設定", "bi bi-gear"),
    new NavigationItem(Page.AppInfo, "アプリ情報", "bi bi-info-circle"),
    new NavigationItem(
      Page.AlertPlayer,
      "プレイヤーリスト",
      "bi bi-person-lines-fill"
    ),
  ];

  export const FUNCS: NavigationItem<Func>[] = [
    new NavigationItem(Func.Reload, "リロード", "bi bi-arrow-clockwise"),
    new NavigationItem(Func.Screenshot, "スクリーンショット", "bi bi-camera"),
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

  export const NATION_ICON = {
    japan: iconJapan,
    usa: iconUsa,
    ussr: iconUssr,
    germany: iconGermany,
    uk: iconUk,
    france: iconFrance,
    italy: iconItaly,
    pan_asia: iconPanAsia,
    europe: iconEurope,
    netherlands: iconNetherlands,
    commonwealth: iconCommonWealth,
    pan_america: iconPanAmerica,
    spain: iconSpain,
    none: iconNone,
  };

  export const SKILL_LEVELS = [
    new SkillLevelItem(SkillLevel.Bad, 11, "cv", 0, 750, 0, 0.6, 0, 47),
    new SkillLevelItem(
      SkillLevel.BelowAvg,
      10,
      "bb",
      750,
      1100,
      0.6,
      0.8,
      47,
      50
    ),
    new SkillLevelItem(SkillLevel.Avg, 9, "bb", 1100, 1350, 0.8, 1.0, 50, 52),
    new SkillLevelItem(SkillLevel.Good, 8, "cl", 1350, 1550, 1.0, 1.2, 52, 54),
    new SkillLevelItem(
      SkillLevel.VeryGood,
      7,
      "cl",
      1550,
      1750,
      1.2,
      1.4,
      54,
      56
    ),
    new SkillLevelItem(SkillLevel.Great, 6, "dd", 1750, 2100, 1.4, 1.5, 56, 60),
    new SkillLevelItem(
      SkillLevel.Unicum,
      5,
      "dd",
      2100,
      2450,
      1.5,
      1.6,
      60,
      65
    ),
    new SkillLevelItem(
      SkillLevel.SuperUnicum,
      4,
      "ss",
      2450,
      9999,
      1.6,
      10,
      65,
      100
    ),
  ];
}
