import { DispName } from "src/lib/DispName";
import type { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { type ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { Damage } from "src/lib/column/model/Damage";
import { PR } from "src/lib/column/model/PR";
import { WinRate } from "src/lib/column/model/WinRate";
import {
  type OptionalBattle,
  type OptionalSummary,
  type StatsCategory,
} from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import { domain } from "wailsjs/go/models";

type SummaryColumn = AbstractColumn<any> & ISummaryColumn;

export interface Summary {
  meta: SummaryMeta;
  values: SummaryValues;
}

interface SummaryMeta {
  headers: SummaryHeader[];
  columnNames: string[];
}

interface SummaryHeader {
  title: string;
  colspan: number;
}

interface SummaryValues {
  friends: string[];
  enemies: string[];
  diffs: SummaryDiff[];
}

interface SummaryDiff {
  diff: string;
  colorCode: string;
}

export namespace Summary {
  export const calculate = (
    battle: OptionalBattle,
    excludedIDs: number[],
    config: domain.UserConfig,
  ): OptionalSummary => {
    if (!validate(battle)) {
      return undefined;
    }

    const { columns, headers } = deriveColumns(config);

    const result: Summary = {
      meta: { headers, columnNames: [] },
      values: { friends: [], enemies: [], diffs: [] },
    };
    columns.forEach((column) => {
      const filtered = battle.teams.map((team) => {
        return team.players.filter(
          (player) =>
            !isExcluded(player, excludedIDs) &&
            isMinBattlesOrMore(player, config, column.getCategory()),
        );
      });

      const [friendMean, enemyMean] = filtered.map((players) =>
        mean(players, column),
      );

      const digit = column.getDigit();

      result.meta.columnNames.push(column.minDisplayName);
      result.values.friends.push(friendMean.toFixed(digit));
      result.values.enemies.push(enemyMean.toFixed(digit));
      result.values.diffs.push(deriveDiff(friendMean, enemyMean, digit));
    });

    return result;
  };
}

const validate = (battle: OptionalBattle): battle is domain.Battle => {
  if (!battle) {
    return false;
  }

  if (battle.teams.length < 2) {
    return false;
  }

  return true;
};

const deriveColumns = (
  config: domain.UserConfig,
): {
  columns: SummaryColumn[];
  headers: SummaryHeader[];
} => {
  const shipCols: SummaryColumn[] = [
    new PR(undefined, config, "ship"),
    new Damage(undefined, config, "ship"),
    new WinRate(undefined, config, "ship"),
  ];

  const overallCols: SummaryColumn[] = [
    new PR(undefined, config, "overall"),
    new Damage(undefined, config, "overall"),
    new WinRate(undefined, config, "overall"),
  ];

  const columns = shipCols.concat(overallCols);
  const headers = [
    {
      title: DispName.COLUMN_CATEGORIES.get("ship")!,
      colspan: shipCols.length,
    },
    {
      title: DispName.COLUMN_CATEGORIES.get("overall")!,
      colspan: overallCols.length,
    },
  ];

  return { columns, headers };
};

const isExcluded = (player: domain.Player, excludedIDs: number[]): boolean => {
  const accountID = player.player_info.id;
  return accountID === 0 || excludedIDs.includes(accountID);
};

const isMinBattlesOrMore = (
  player: domain.Player,
  config: domain.UserConfig,
  category: StatsCategory,
): boolean => {
  const battles = toPlayerStats(player, config.stats_pattern)[category].battles;
  const teamAverage = config.team_average;

  let minBattles: number;
  switch (category) {
    case "ship":
      minBattles = teamAverage.min_ship_battles;
      break;
    case "overall":
      minBattles = teamAverage.min_overall_battles;
      break;
  }

  return battles >= minBattles;
};

const mean = (players: domain.Player[], column: SummaryColumn): number => {
  const values = players.map((player) => column.getValue(player));
  return values.length !== 0
    ? values.reduce((a, b) => a + b, 0) / values.length
    : 0;
};

const deriveDiff = (
  friend: number,
  enemy: number,
  digit: number,
): SummaryDiff => {
  const diffNum = friend - enemy;
  const sign = diffNum > 0 ? "+" : "";

  let colorCode = "";
  switch (true) {
    case diffNum > 0:
      colorCode = "#99d02b";
      break;
    case diffNum < 0:
      colorCode = "#fc4e32";
      break;
  }
  const diff = sign + diffNum.toFixed(digit);

  return { diff, colorCode };
};
