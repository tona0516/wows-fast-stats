import type { vo } from "wailsjs/go/models";
import Const from "./Const";

const SHIP_STATS_KEY = "ship_stats";
const OVERALL_STATS_KEY = "overall_stats";

const keys = [
  { key1: SHIP_STATS_KEY, key2: "pr" },
  { key1: SHIP_STATS_KEY, key2: "damage" },
  { key1: SHIP_STATS_KEY, key2: "win_rate" },
  { key1: SHIP_STATS_KEY, key2: "kd_rate" },
  { key1: OVERALL_STATS_KEY, key2: "damage" },
  { key1: OVERALL_STATS_KEY, key2: "win_rate" },
  { key1: OVERALL_STATS_KEY, key2: "kd_rate" },
];

export class AverageFactor {
  shipStatsCount: number;
  overallCount: number;
  labels: string[];
  friends: string[];
  enemies: string[];
  diffs: { value: string; colorClass: string }[];
}

export class Average {
  private battle: vo.Battle;

  constructor(battle: vo.Battle) {
    this.battle = battle;
  }

  calc(excludePlayerIDs: number[]): AverageFactor {
    const targetFriends = this.battle.teams[0].players.filter(
      (it) => !excludePlayerIDs.includes(it.player_info.id)
    );
    const targetEnemies = this.battle.teams[1].players.filter(
      (it) => !excludePlayerIDs.includes(it.player_info.id)
    );

    let shipStatsCount = 0;
    let overallCount = 0;
    const labels: string[] = [];
    const friends: string[] = [];
    const diffs: { value: string; colorClass: string }[] = [];
    const enemies: string[] = [];
    keys.forEach((it) => {
      const friendAvg = average(it.key1, it.key2, targetFriends);
      const enemyAvg = average(it.key1, it.key2, targetEnemies);
      const diff = friendAvg - enemyAvg;
      let sign = diff > 0 ? "+" : "";
      let colorClass = "";
      if (diff > 0) {
        colorClass = "higher";
      } else if (diff < 0) {
        colorClass = "lower";
      }

      labels.push(Const.COLUMN_NAMES[it.key2].min);
      friends.push(friendAvg.toFixed(Const.DIGITS[it.key2]));
      enemies.push(enemyAvg.toFixed(Const.DIGITS[it.key2]));
      diffs.push({
        value: sign + diff.toFixed(Const.DIGITS[it.key2]),
        colorClass: colorClass,
      });

      if (it.key1 === SHIP_STATS_KEY) {
        shipStatsCount++;
      }
      if (it.key1 === OVERALL_STATS_KEY) {
        overallCount++;
      }
    });

    return {
      shipStatsCount: shipStatsCount,
      overallCount: overallCount,
      labels: labels,
      friends: friends,
      enemies: enemies,
      diffs: diffs,
    };
  }
}

function average(key1: string, key2: string, players: vo.Player[]): number {
  const values = players
    .filter((it) => it[key1]["battles"] !== 0)
    .map((it) => it[key1][key2] as number);

  if (values.length === 0) {
    return 0;
  }

  return values.reduce((a, b) => a + b, 0) / values.length;
}
