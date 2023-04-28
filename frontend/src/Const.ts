import type { vo } from "wailsjs/go/models.js";
import type { Column } from "./Column";

namespace Const {
    export const DEFAULT_USER_CONFIG: vo.UserConfig = {
        install_path: "",
        appid: "",
        font_size: "medium",
        displays: {
            basic: {
                player_name: true,
                ship_info: true
            },
            ship: {
                pr: true,
                damage: true,
                win_rate: true,
                kd_rate: false,
                win_survived_rate: false,
                lose_survived_rate: false,
                exp: false,
                battles: true
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
                using_tier_rate: false
            },
            convertValues: function (a: any, classs: any, asMap?: boolean) {
                throw new Error("Function not implemented.");
            }
        },
        save_screenshot: false,
        save_temp_arena_info: false,
        convertValues: function (a: any, classs: any, asMap?: boolean) {
            throw new Error("Function not implemented.");
        }
    };

    export const COLUMN_NAMES: {[key: string]: Column} = {
        basic: {minName: "基本情報", fullName: "基本情報"},
        ship: {minName: "艦", fullName: "艦成績"},
        overall: {minName: "総合", fullName: "総合成績"},
        player_name: {minName: "プレイヤー", fullName: "プレイヤー"},
        ship_info: {minName: "艦", fullName: "艦情報"},
        pr: {minName: "PR", fullName: "Personal Rating"},
        damage: {minName: "Dmg", fullName: "ダメージ"},
        win_rate: {minName: "勝率", fullName: "勝率"},
        kd_rate: {minName: "K/D", fullName: "K/D比"},
        win_survived_rate: {minName: "生存率(勝)", fullName: "勝利生存率"},
        lose_survived_rate: {minName: "生存率(負)", fullName: "敗北生存率"},
        exp: {minName: "Exp", fullName: "経験値"},
        battles: {minName: "戦闘数", fullName: "戦闘数"},
        avg_tier: {minName: "平均T", fullName: "平均Tier"},
        using_ship_type_rate: {minName: "艦割合", fullName: "艦種別プレイ割合"},
        using_tier_rate: {minName: "T割合", fullName: "ティア別プレイ割合"},
    }

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
    }
}

export default Const
