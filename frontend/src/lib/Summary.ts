import { Damage } from "src/lib/column/Damage";
import { KDRate } from "src/lib/column/KDRate";
import { PR } from "src/lib/column/PR";
import { WinRate } from "src/lib/column/WinRate";
import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { type ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { type OptionalBattle, type OptionalSummary } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import { domain } from "wailsjs/go/models";

export class Summary {
  constructor(
    readonly shipColspan: number,
    readonly overallColspan: number,
    readonly labels: string[],
    readonly friends: string[],
    readonly enemies: string[],
    readonly diffs: { value: string; colorClass: string }[],
  ) {}

  static calculate(
    battle: OptionalBattle,
    excludes: number[],
    userConfig: domain.UserConfig,
  ): OptionalSummary {
    if (!battle) {
      return undefined;
    }

    const columns: (AbstractSingleColumn<any> & ISummaryColumn)[] = [
      new PR(userConfig, "ship"),
      new Damage(userConfig, "ship"),
      new WinRate(userConfig, "ship"),
      new KDRate(userConfig, "ship"),
      new Damage(userConfig, "overall"),
      new WinRate(userConfig, "overall"),
      new KDRate(userConfig, "overall"),
    ];

    const teams = [battle.teams[0], battle.teams[1]];

    const labels: string[] = [];
    const friends: string[] = [];
    const enemies: string[] = [];
    const diffs: { value: string; colorClass: string }[] = [];

    columns.forEach((column) => {
      const [filteredFriends, filteredEnemies] = teams.map((team) => {
        return team.players.filter((player) => {
          const battles = toPlayerStats(player, userConfig.stats_pattern)[
            column.category()
          ].battles;
          const teamAverage = userConfig.team_average;

          let minBattles: number;
          switch (column.category()) {
            case "ship":
              minBattles = teamAverage.min_ship_battles;
            case "overall":
              minBattles = teamAverage.min_overall_battles;
          }

          const accountID = player.player_info.id;
          return (
            accountID !== 0 &&
            !excludes.includes(accountID) &&
            battles >= minBattles
          );
        });
      });

      const [friendMean, enemyMean] = [filteredFriends, filteredEnemies].map(
        (team) => {
          const values = team.map((player) => column.value(player));
          return values.length !== 0
            ? values.reduce((a, b) => a + b, 0) / values.length
            : 0;
        },
      );

      const diff = friendMean - enemyMean;
      const sign = diff > 0 ? "+" : "";
      const colorClass = diff > 0 ? "higher" : "lower";
      const digit = column.digit();

      labels.push(column.minDisplayName());
      friends.push(friendMean.toFixed(digit));
      enemies.push(enemyMean.toFixed(digit));
      diffs.push({ value: sign + diff.toFixed(digit), colorClass: colorClass });
    });

    return {
      shipColspan: columns.filter((it) => it.category() === "ship").length,
      overallColspan: columns.filter((it) => it.category() === "overall")
        .length,
      labels: labels,
      friends: friends,
      enemies: enemies,
      diffs: diffs,
    };
  }
}
