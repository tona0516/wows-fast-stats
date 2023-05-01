import type { vo } from "wailsjs/go/models";
import Const from "./Const";

export class AverageFactor {
    label: string
    friend: string
    enemy: string
    diff: string
    colorClass: string
}

export class Average {
    private battle: vo.Battle

    constructor(battle: vo.Battle) {
        this.battle = battle
    }

    calc(excludePlayerID: number[]): AverageFactor[] {
        const friends = this.battle.teams[0].players.filter((it) => !excludePlayerID.includes(it.player_info.id))
        const enemies = this.battle.teams[1].players.filter((it) => !excludePlayerID.includes(it.player_info.id))

        const result: AverageFactor[] = [];
        [
            {key1: "ship_stats", key2: "pr"},
            {key1: "ship_stats", key2: "damage"},
            {key1: "ship_stats", key2: "win_rate"},
            {key1: "ship_stats", key2: "kd_rate"},
            {key1: "overall_stats", key2: "damage"},
            {key1: "overall_stats", key2: "win_rate"},
            {key1: "overall_stats", key2: "kd_rate"},
        ].forEach(it => {
            result.push(calcFactor(it.key1, it.key2, friends, enemies))
        });

        return result
    }
}

function calcFactor(key1: string, key2: string, friends: vo.Player[], enemies: vo.Player[]): AverageFactor {
    const friend = friends.map((it) => it[key1][key2]).reduce((a, b) => a + b, 0) / friends.length || 0
    const enemy = enemies.map((it) => it[key1][key2]).reduce((a, b) => a + b, 0) / enemies.length || 0
    const diff = friend - enemy
    let colorClass = "";
    let sign = "";
    if (diff > 0) {
        sign = "+";
        colorClass = "higher";
    } else if (diff < 0) {
        colorClass = "lower";
    }

    return {
        label: Const.COLUMN_NAMES[key1].min + ":" + Const.COLUMN_NAMES[key2].min,
        friend: friend.toFixed(Const.DIGITS[key2]),
        enemy: enemy.toFixed(Const.DIGITS[key2]),
        diff: sign + diff.toFixed(Const.DIGITS[key2]),
        colorClass: colorClass,
    }
}
