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
    readonly tableInfo: {
      shipColspan: number;
      overallColspan: number;
    },
    readonly values: {
      label: string;
      friend: string;
      enemy: string;
      diff: { value: string; color: string };
    }[],
  ) {}

  static calculate(
    battle: OptionalBattle,
    excludedIDs: number[],
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
    const shipColspan = columns.filter((it) => it.category() === "ship").length;
    const overallColspan = columns.filter(
      (it) => it.category() === "overall",
    ).length;

    if (battle.teams.length < 2) {
      return undefined;
    }

    const teams = [battle.teams[0], battle.teams[1]];

    let values: {
      label: string;
      friend: string;
      enemy: string;
      diff: { value: string; color: string };
    }[] = [];

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
            !excludedIDs.includes(accountID) &&
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

      const diffNum = friendMean - enemyMean;
      let sign = diffNum > 0 ? "+" : "";
      let colorCode = "";
      if (diffNum > 0) {
        colorCode = "#99d02b";
      } else if (diffNum < 0) {
        colorCode = "#fc4e32";
      }
      const diffStr = sign + diffNum.toFixed(digit);

      values.push({
        label: column.minDisplayName(),
        friend: friendMean.toFixed(digit),
        enemy: enemyMean.toFixed(digit),
        diff: { value: diffStr, color: colorCode },
      });
    });

    return {
      tableInfo: {
        shipColspan: shipColspan,
        overallColspan: overallColspan,
      },
      values: values,
    };
  }
}
