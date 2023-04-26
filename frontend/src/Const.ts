import type { vo } from "wailsjs/go/models.js";

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
}

export default Const
