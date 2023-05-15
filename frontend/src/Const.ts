import type { vo } from "wailsjs/go/models.js";

namespace Const {
  export const BASE_NUMBERS_URL = "https://asia.wows-numbers.com/";

  export const DEFAULT_USER_CONFIG: vo.UserConfig = {
    install_path: "",
    appid: "",
    font_size: "medium",
    displays: {
      basic: {
        is_in_avg: true,
        player_name: true,
        ship_info: true,
      },
      ship: {
        pr: true,
        damage: true,
        win_rate: true,
        kd_rate: false,
        exp: false,
        battles: true,
        survived_rate: false,
        hit_rate: false,
      },
      overall: {
        damage: true,
        win_rate: true,
        kd_rate: false,
        exp: false,
        battles: true,
        survived_rate: false,
        avg_tier: false,
        using_ship_type_rate: false,
        using_tier_rate: false,
      },
      convertValues: function (_a: any, _classs: any, _asMap?: boolean) {
        throw new Error("Function not implemented.");
      },
    },
    save_screenshot: false,
    save_temp_arena_info: false,
    convertValues: function (_a: any, _classs: any, _asMap?: boolean) {
      throw new Error("Function not implemented.");
    },
  };

  export const COLUMN_NAMES = {
    basic: { min: "基本情報", full: "基本情報" },
    ship_stats: { min: "艦", full: "艦成績" },
    overall_stats: { min: "総合", full: "総合成績" },
    is_in_avg: { min: "", full: "" },
    player_name: { min: "プレイヤー", full: "プレイヤー" },
    ship_info: { min: "艦", full: "艦情報" },
    pr: { min: "PR", full: "Personal Rating" },
    damage: { min: "Dmg", full: "ダメージ" },
    win_rate: { min: "勝率", full: "勝率" },
    kd_rate: { min: "K/D", full: "K/D比" },
    survived_rate: { min: "生存率(勝|負)", full: "生存率" },
    exp: { min: "Exp", full: "経験値" },
    battles: { min: "戦闘数", full: "戦闘数" },
    avg_tier: { min: "平均T", full: "平均Tier" },
    using_ship_type_rate: { min: "艦割合", full: "艦種別プレイ割合" },
    using_tier_rate: { min: "T割合", full: "ティア別プレイ割合" },
    hit_rate: { min: "Hit率(主|魚)", full: "命中率" },
  };

  export const DIGITS: { [key: string]: number } = {
    pr: 0,
    damage: 0,
    win_rate: 1,
    kd_rate: 1,
    survived_rate: 1,
    hit_rate: 1,
    exp: 0,
    battles: 0,
    avg_tier: 1,
    ship_type_rate: 1,
    tier_rate: 1,
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
    "": "#00000000",
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
    "": "#00000000",
    bad: "#ff382d",
    belowAvg: "#fd9234",
    avg: "#ffd351",
    good: "#57e500",
    veryGood: "#44b200",
    great: "#02f7da",
    unicum: "#da6ff5",
    superUnicum: "#bf15ee",
  };
}

export default Const;
