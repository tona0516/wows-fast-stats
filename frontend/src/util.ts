import { format, fromUnixTime } from "date-fns";
import type { domain } from "../wailsjs/go/models";
import { Const } from "./Const";
import { DisplayPattern, StatsCategory } from "./enums";
import { SkillLevelConverter } from "./RankConverter";

export function colors(
  key: string,
  value: number,
  player: domain.Player,
  statsCategory: StatsCategory,
  skillColor: domain.SkillColor
): string {
  switch (key) {
    case "pr":
      return SkillLevelConverter.fromPR(value, skillColor).toTextColorCode();
    case "damage":
      if (statsCategory === StatsCategory.Ship) {
        return SkillLevelConverter.fromDamage(
          value,
          player.ship_info.avg_damage,
          skillColor
        ).toTextColorCode();
      }
      return "";
    case "win_rate":
      return SkillLevelConverter.fromWinRate(
        value,
        skillColor
      ).toTextColorCode();
    default:
      return "";
  }
}

export function values(
  player: domain.Player,
  statsPattern: string,
  statsCategory: StatsCategory,
  key: string
): any {
  return player[statsPattern][statsCategory][key];
}

export interface SummaryResult {
  shipColspan: number;
  overallColspan: number;
  labels: string[];
  friends: string[];
  enemies: string[];
  diffs: { value: string; colorClass: string }[];
}

export function summary(
  battle: domain.Battle,
  excludes: number[],
  userConfig: domain.UserConfig
): SummaryResult {
  if (!battle) {
    return undefined;
  }

  const items: { category: StatsCategory; key: string }[] = [
    { category: StatsCategory.Ship, key: "pr" },
    { category: StatsCategory.Ship, key: "damage" },
    { category: StatsCategory.Ship, key: "win_rate" },
    { category: StatsCategory.Ship, key: "kd_rate" },
    { category: StatsCategory.Overall, key: "damage" },
    { category: StatsCategory.Overall, key: "win_rate" },
    { category: StatsCategory.Overall, key: "kd_rate" },
  ];

  const [shipColspan, overallColspan] = [
    StatsCategory.Ship,
    StatsCategory.Overall,
  ].map((category) => {
    return items.filter((it) => it.category === category).length;
  });

  const labels: string[] = [];
  const friends: string[] = [];
  const enemies: string[] = [];
  const diffs: { value: string; colorClass: string }[] = [];
  items.forEach((it) => {
    const [filteredFriends, filteredEnemies] = [0, 1].map((i) => {
      let minBattle = 1;
      if (it.category === "ship") {
        minBattle = userConfig.team_average.min_ship_battles;
      } else if (it.category === "overall") {
        minBattle = userConfig.team_average.min_overall_battles;
      }

      return battle.teams[i].players.filter(
        (p) =>
          p.player_info.id !== 0 &&
          !excludes.includes(p.player_info.id) &&
          values(p, userConfig.stats_pattern, it.category, "battles") >=
            minBattle
      );
    });

    const [friendMean, enemyMean] = [filteredFriends, filteredEnemies].map(
      (team) => {
        return mean(team, it.category, userConfig.stats_pattern, it.key);
      }
    );

    const diff = friendMean - enemyMean;
    let sign = diff > 0 ? "+" : "";
    let colorClass = "";
    if (diff > 0) {
      colorClass = "higher";
    } else if (diff < 0) {
      colorClass = "lower";
    }

    const digit = userConfig.custom_digit[it.key];

    labels.push(Const.COLUMN_NAMES[it.key].min);
    friends.push(friendMean.toFixed(digit));
    enemies.push(enemyMean.toFixed(digit));
    diffs.push({
      value: sign + diff.toFixed(digit),
      colorClass: colorClass,
    });
  });

  return {
    shipColspan: shipColspan,
    overallColspan: overallColspan,
    labels: labels,
    friends: friends,
    enemies: enemies,
    diffs: diffs,
  };
}

export function clanURL(player: domain.Player): string {
  return (
    Const.BASE_NUMBERS_URL +
    "clan/" +
    player.player_info.clan.id +
    "," +
    player.player_info.clan.tag
  );
}

export function playerURL(player: domain.Player): string {
  return (
    Const.BASE_NUMBERS_URL +
    "player/" +
    player.player_info.id +
    "," +
    player.player_info.name
  );
}

export function shipURL(player: domain.Player): string {
  return (
    Const.BASE_NUMBERS_URL +
    "ship/" +
    player.ship_info.id +
    "," +
    player.ship_info.name.replaceAll(" ", "-")
  );
}

export function tierString(value: number): string {
  if (value === 11) return "★";

  const decimal = [10, 9, 5, 4, 1];
  const romanNumeral = ["X", "IX", "V", "IV", "I"];

  let romanized = "";

  for (var i = 0; i < decimal.length; i++) {
    while (decimal[i] <= value) {
      romanized += romanNumeral[i];
      value -= decimal[i];
    }
  }

  return romanized;
}

export function decideDisplayPattern(
  player: domain.Player,
  statsPattern: string
): DisplayPattern {
  if (player.player_info.is_hidden) {
    return DisplayPattern.Private;
  }

  if (
    player.player_info.id === 0 ||
    player[statsPattern].overall.battles === 0
  ) {
    return DisplayPattern.NoData;
  }

  if (player[statsPattern].ship.battles === 0) {
    return DisplayPattern.NoShipStats;
  }

  return DisplayPattern.Full;
}

export function toDateForDisplay(unixtime: number): string {
  return format(fromUnixTime(unixtime), "yyyy/MM/dd HH:mm:ss");
}

export function toDateForFilename(unixtime: number): string {
  return format(fromUnixTime(unixtime), "yyyy-MM-dd-HH-mm-ss");
}

function mean(
  players: domain.Player[],
  statsCategory: StatsCategory,
  statsPattern: string,
  key: string
): number {
  let values: number[] = [];

  // Note: PR is -1 when expected values can't retrieve.
  if (key == "pr") {
    values = players
      .filter((it) => it[statsPattern][statsCategory][key] !== -1)
      .map((it) => it[statsPattern][statsCategory][key]);
  } else {
    values = players.map((it) => it[statsPattern][statsCategory][key]);
  }

  if (values.length === 0) {
    return 0;
  }

  return values.reduce((a, b) => a + b, 0) / values.length;
}

export function toTierGroup(tier: number): string {
  if (tier >= 1 && tier <= 4) {
    return "low";
  }
  if (tier >= 5 && tier <= 7) {
    return "middle";
  }
  if (tier >= 8) {
    return "high";
  }
  return "";
}
