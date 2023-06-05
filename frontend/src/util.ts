import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import type { StatsPattern } from "./StatsPattern";
import Const from "./Const";
import type { StatsCategory } from "./StatsCategory";

export function values(
  player: vo.Player,
  displayPattern: DisplayPattern,
  statsPattern: StatsPattern,
  statsCategory: StatsCategory,
  key: string
): any {
  if (statsCategory === "ship") {
    if (displayPattern === "full") {
      if (statsPattern === "pvp_solo") {
        return player.ship_stats_solo[key];
      }
      if (statsPattern === "pvp_all") {
        return player.ship_stats[key];
      }
    }

    return undefined;
  }

  if (statsCategory === "overall") {
    if (["full", "noshipstats"].includes(displayPattern)) {
      if (statsPattern === "pvp_solo") {
        return player.overall_stats_solo[key];
      }
      if (statsPattern === "pvp_all") {
        return player.overall_stats[key];
      }
    }

    return undefined;
  }

  return undefined;
}

export class SummaryResult {
  shipStatsCount: number;
  overallStatsCount: number;
  labels: string[];
  friends: string[];
  enemies: string[];
  diffs: { value: string; colorClass: string }[];
}

export function summary(
  battle: vo.Battle,
  excludes: number[],
  statsPattern: StatsPattern
): SummaryResult {
  const filteredFriends = battle.teams[0].players.filter(
    (it) => !excludes.includes(it.player_info.id)
  );
  const filteredEnemies = battle.teams[1].players.filter(
    (it) => !excludes.includes(it.player_info.id)
  );

  const keys = [
    { statsKey: toStatsKey("ship", statsPattern), valueKey: "pr" },
    { statsKey: toStatsKey("ship", statsPattern), valueKey: "damage" },
    { statsKey: toStatsKey("ship", statsPattern), valueKey: "win_rate" },
    { statsKey: toStatsKey("ship", statsPattern), valueKey: "kd_rate" },
    { statsKey: toStatsKey("overall", statsPattern), valueKey: "damage" },
    { statsKey: toStatsKey("overall", statsPattern), valueKey: "win_rate" },
    { statsKey: toStatsKey("overall", statsPattern), valueKey: "kd_rate" },
  ];

  const labels: string[] = [];
  const friends: string[] = [];
  const enemies: string[] = [];
  const diffs: { value: string; colorClass: string }[] = [];
  keys.forEach((it) => {
    const friendMean = mean(filteredFriends, it.statsKey, it.valueKey);
    const enemyMean = mean(filteredEnemies, it.statsKey, it.valueKey);

    const diff = friendMean - enemyMean;
    let sign = diff > 0 ? "+" : "";
    let colorClass = "";
    if (diff > 0) {
      colorClass = "higher";
    } else if (diff < 0) {
      colorClass = "lower";
    }

    const digit = Const.DIGITS[it.valueKey];

    labels.push(Const.COLUMN_NAMES[it.valueKey].min);
    friends.push(friendMean.toFixed(digit));
    enemies.push(enemyMean.toFixed(digit));
    diffs.push({
      value: sign + diff.toFixed(digit),
      colorClass: colorClass,
    });
  });

  return {
    shipStatsCount: keys.filter((it) => it.statsKey.startsWith("ship")).length,
    overallStatsCount: keys.filter((it) => it.statsKey.startsWith("overall"))
      .length,
    labels: labels,
    friends: friends,
    enemies: enemies,
    diffs: diffs,
  };
}

// TODO Refactoring
function toStatsKey(
  statsCategory: StatsCategory,
  statsPattern: StatsPattern
): string {
  if (statsCategory === "ship" && statsPattern === "pvp_all") {
    return "ship_stats";
  }
  if (statsCategory === "ship" && statsPattern === "pvp_solo") {
    return "ship_stats_solo";
  }
  if (statsCategory === "overall" && statsPattern === "pvp_all") {
    return "overall_stats";
  }
  if (statsCategory === "overall" && statsPattern === "pvp_solo") {
    return "overall_stats_solo";
  }

  return undefined;
}

function mean(
  players: vo.Player[],
  statsKey: string,
  valueKey: string
): number {
  let values: number[] = [];
  // Note: PR is -1 when expected values can't retrieve.
  if (valueKey == "pr") {
    values = players
      .filter(
        (it) => it[statsKey]["battles"] !== 0 && it[statsKey][valueKey] >= 0
      )
      .map((it) => it[statsKey][valueKey] as number);
  } else {
    values = players
      .filter((it) => it[statsKey]["battles"] !== 0)
      .map((it) => it[statsKey][valueKey] as number);
  }

  if (values.length === 0) {
    return 0;
  }

  return values.reduce((a, b) => a + b, 0) / values.length;
}
