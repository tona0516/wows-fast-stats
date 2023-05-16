import type { vo } from "wailsjs/go/models";
import Const from "./Const";

const SHIP_STATS_KEY = "ship_stats";
const OVERALL_STATS_KEY = "overall_stats";

const keys = [
  { category: SHIP_STATS_KEY, name: "pr" },
  { category: SHIP_STATS_KEY, name: "damage" },
  { category: SHIP_STATS_KEY, name: "win_rate" },
  { category: SHIP_STATS_KEY, name: "kd_rate" },
  { category: OVERALL_STATS_KEY, name: "damage" },
  { category: OVERALL_STATS_KEY, name: "win_rate" },
  { category: OVERALL_STATS_KEY, name: "kd_rate" },
];

function mean(category: string, name: string, players: vo.Player[]): number {
  let values: number[] = [];

  // Note: PR is -1 when expected values can't retrieve.
  if (name == "pr") {
    values = players
      .filter((it) => it[category]["battles"] !== 0 && it[category][name] >= 0)
      .map((it) => it[category][name] as number);
  } else {
    values = players
      .filter((it) => it[category]["battles"] !== 0)
      .map((it) => it[category][name] as number);
  }

  if (values.length === 0) {
    return 0;
  }

  return values.reduce((a, b) => a + b, 0) / values.length;
}

export class SummaryResult {
  shipStatsCount: number;
  overallStatsCount: number;
  labels: string[];
  friends: string[];
  enemies: string[];
  diffs: { value: string; colorClass: string }[];
}

export class Summary {
  private battle: vo.Battle;

  constructor(battle: vo.Battle) {
    this.battle = battle;
  }

  calc(excludePlayerIDs: number[]): SummaryResult {
    const targetFriends = this.battle.teams[0].players.filter(
      (it) => !excludePlayerIDs.includes(it.player_info.id)
    );
    const targetEnemies = this.battle.teams[1].players.filter(
      (it) => !excludePlayerIDs.includes(it.player_info.id)
    );

    const labels: string[] = [];
    const friends: string[] = [];
    const enemies: string[] = [];
    const diffs: { value: string; colorClass: string }[] = [];
    keys.forEach((it) => {
      const friendMean = mean(it.category, it.name, targetFriends);
      const enemyMean = mean(it.category, it.name, targetEnemies);

      const diff = friendMean - enemyMean;
      let sign = diff > 0 ? "+" : "";
      let colorClass = "";
      if (diff > 0) {
        colorClass = "higher";
      } else if (diff < 0) {
        colorClass = "lower";
      }

      const digit = Const.DIGITS[it.name];

      labels.push(Const.COLUMN_NAMES[it.name].min);
      friends.push(friendMean.toFixed(digit));
      enemies.push(enemyMean.toFixed(digit));
      diffs.push({
        value: sign + diff.toFixed(digit),
        colorClass: colorClass,
      });
    });

    return {
      shipStatsCount: keys.filter((it) => it.category === SHIP_STATS_KEY)
        .length,
      overallStatsCount: keys.filter((it) => it.category === OVERALL_STATS_KEY)
        .length,
      labels: labels,
      friends: friends,
      enemies: enemies,
      diffs: diffs,
    };
  }
}
