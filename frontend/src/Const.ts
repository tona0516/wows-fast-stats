import type { vo } from "wailsjs/go/models.js";

namespace Const {
    export const DEFAULT_USER_CONFIG: vo.UserConfig = {
        install_path: "",
        appid: "",
        font_size: "medium",
        displays: {
            pr: true,
            ship_damage: true,
            ship_win_rate: true,
            ship_kd_rate: true,
            ship_win_survived_rate: false,
            ship_lose_survived_rate: false,
            ship_exp: false,
            ship_battles: true,
            player_damage: true,
            player_win_rate: true,
            player_kd_rate: true,
            player_win_survived_rate: false,
            player_lose_survived_rate: false,
            player_exp: false,
            player_battles: true,
            player_avg_tier: false,
        },
        convertValues: function (a: any, classs: any, asMap?: boolean) {
            throw new Error("Function not implemented.");
        }
    };
}

export default Const
