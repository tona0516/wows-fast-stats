import { Damage } from "src/lib/column/Damage";
import { PR } from "src/lib/column/PR";
import { WinRate } from "src/lib/column/WinRate";
import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { type ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { type OptionalBattle, type OptionalSummary } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import { domain } from "wailsjs/go/models";

export class Summary {
  constructor(
    readonly friends: { [label: string]: number },
    readonly enemies: { [label: string]: number },
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
      new PR(userConfig, "overall"),
      new Damage(userConfig, "overall"),
      new WinRate(userConfig, "overall"),
    ];

    const teams = [battle.teams[0], battle.teams[1]];

    const friends: { [label: string]: number } = {};
    const enemies: { [label: string]: number } = {};

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

      const digit = column.digit();

      const label = `${column.category()}:${column.minDisplayName()}`;

      friends[label] = Number(friendMean.toFixed(digit));
      enemies[label] = Number(enemyMean.toFixed(digit));
    });

    return {
      friends: friends,
      enemies: enemies,
    };
  }
}
