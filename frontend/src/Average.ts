import type { vo } from "wailsjs/go/models";
import Const from "./Const";

const keys = [
    {key1: "ship_stats", key2: "pr"},
    {key1: "ship_stats", key2: "damage"},
    {key1: "ship_stats", key2: "win_rate"},
    {key1: "ship_stats", key2: "kd_rate"},
    {key1: "overall_stats", key2: "damage"},
    {key1: "overall_stats", key2: "win_rate"},
    {key1: "overall_stats", key2: "kd_rate"},
]

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

    calc(excludePlayerIDs: number[]): AverageFactor[] {
        const friends = this.battle.teams[0].players.filter((it) => !excludePlayerIDs.includes(it.player_info.id))
        const enemies = this.battle.teams[1].players.filter((it) => !excludePlayerIDs.includes(it.player_info.id))
        return keys.map((it) => calcFactor(it.key1, it.key2, friends, enemies))
    }
}

function calcFactor(key1: string, key2: string, friends: vo.Player[], enemies: vo.Player[]): AverageFactor {
    const friendAvg = average(key1, key2, friends)
    const enemyAvg = average(key1, key2, enemies)
    const diff = friendAvg - enemyAvg

    let sign = diff > 0 ? "+" : "";
    let colorClass = "";
    if (diff > 0) {
        colorClass = "higher";
    } else if (diff < 0) {
        colorClass = "lower";
    }

    return {
        label: Const.COLUMN_NAMES[key1].min + ":" + Const.COLUMN_NAMES[key2].min,
        friend: friendAvg.toFixed(Const.DIGITS[key2]),
        enemy: enemyAvg.toFixed(Const.DIGITS[key2]),
        diff: sign + diff.toFixed(Const.DIGITS[key2]),
        colorClass: colorClass,
    }
}

function average(key1: string, key2: string, players: vo.Player[]): number {
    const values = players
        .filter((it) => it[key1]["battles"] !== 0)
        .map((it) => it[key1][key2] as number)

    if (values.length === 0) {
        return 0
    }

    return values.reduce((a, b) => a + b, 0) / values.length
}
