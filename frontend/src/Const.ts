import type { vo } from "wailsjs/go/models.js";
import type { Column } from "./Column";

namespace Const {
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
        win_survived_rate: false,
        lose_survived_rate: false,
        exp: false,
        battles: true,
      },
      overall: {
        damage: true,
        win_rate: true,
        kd_rate: false,
        win_survived_rate: false,
        lose_survived_rate: false,
        exp: false,
        battles: true,
        avg_tier: false,
        using_ship_type_rate: false,
        using_tier_rate: false,
      },
      convertValues: function (a: any, classs: any, asMap?: boolean) {
        throw new Error("Function not implemented.");
      },
    },
    save_screenshot: false,
    save_temp_arena_info: false,
    convertValues: function (a: any, classs: any, asMap?: boolean) {
      throw new Error("Function not implemented.");
    },
  };

  export const COLUMN_NAMES: { [key: string]: Column } = {
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
    win_survived_rate: { min: "生存率(勝)", full: "勝利生存率" },
    lose_survived_rate: { min: "生存率(負)", full: "敗北生存率" },
    exp: { min: "Exp", full: "経験値" },
    battles: { min: "戦闘数", full: "戦闘数" },
    avg_tier: { min: "平均T", full: "平均Tier" },
    using_ship_type_rate: { min: "艦割合", full: "艦種別プレイ割合" },
    using_tier_rate: { min: "T割合", full: "ティア別プレイ割合" },
  };

  export const DIGITS: { [key: string]: number } = {
    pr: 0,
    damage: 0,
    win_rate: 1,
    kd_rate: 1,
    win_survived_rate: 1,
    lose_survived_rate: 1,
    exp: 0,
    battles: 0,
    avg_tier: 1,
  };
}

export default Const;
