import { ArrayMap } from "src/lib/ArrayMap";
import { DispName } from "src/lib/DispName";
import type { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { type ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { Damage } from "src/lib/column/model/Damage";
import { PR } from "src/lib/column/model/PR";
import { WinRate } from "src/lib/column/model/WinRate";
import {
  type OptionalBattle,
  type OptionalSummary,
  type ShipType,
  type StatsCategory,
} from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import { model } from "wailsjs/go/models";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type SummaryColumn = AbstractStatsColumn<any> & ISummaryColumn;
export type SummaryShipType = ShipType | "all";

export interface Summary {
  meta: SummaryMeta;
  values: ArrayMap<SummaryShipType, SummaryValues>;
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
    config: model.UserConfig,
  ): OptionalSummary => {
    if (!validate(battle)) {
      return undefined;
    }

    const { columns, headers } = deriveColumns(config);
    const result: Summary = {
      meta: { headers, columnNames: [] },
      values: new ArrayMap([
        ["all", { friends: [], enemies: [], diffs: [] }],
        ["cv", { friends: [], enemies: [], diffs: [] }],
        ["bb", { friends: [], enemies: [], diffs: [] }],
        ["cl", { friends: [], enemies: [], diffs: [] }],
        ["dd", { friends: [], enemies: [], diffs: [] }],
        ["ss", { friends: [], enemies: [], diffs: [] }],
      ]),
    };

    columns.forEach((column) => {
      result.meta.columnNames.push(column.header);

      const filtered = battle.teams.map((team) => {
        return team.players.filter(
          (player) =>
            !isExcluded(player, excludedIDs) &&
            isMinBattlesOrMore(player, config, column.category),
        );
      });

      [...result.values.keys()].forEach((shipType) => {
        let origin: model.Player[][];
        if (shipType.toString() === "all") {
          origin = filtered;
        } else {
          origin = filtered.map((it) =>
            it.filter((it) => it.ship_info.type === shipType.toString()),
          );
        }

        const [friendMean, enemyMean] = origin.map((players) =>
          mean(players, column),
        );

        const digit = column.digit();
        result.values.get(shipType)!.friends.push(friendMean.toFixed(digit));
        result.values.get(shipType)!.enemies.push(enemyMean.toFixed(digit));
        result.values
          .get(shipType)!
          .diffs.push(deriveDiff(friendMean, enemyMean, digit));
      });
    });

    return result;
  };
}

const validate = (battle: OptionalBattle): battle is model.Battle => {
  if (!battle) {
    return false;
  }

  if (battle.teams.length < 2) {
    return false;
  }

  return true;
};

const deriveColumns = (
  config: model.UserConfig,
): {
  columns: SummaryColumn[];
  headers: SummaryHeader[];
} => {
  const shipCols: SummaryColumn[] = [
    new PR(config, "ship"),
    new Damage(config, "ship"),
    new WinRate(config, "ship"),
  ];

  const overallCols: SummaryColumn[] = [
    new PR(config, "overall"),
    new Damage(config, "overall"),
    new WinRate(config, "overall"),
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

const isExcluded = (player: model.Player, excludedIDs: number[]): boolean => {
  const accountID = player.player_info.id;
  return accountID === 0 || excludedIDs.includes(accountID);
};

const isMinBattlesOrMore = (
  player: model.Player,
  config: model.UserConfig,
  category: StatsCategory,
): boolean => {
  const battles = toPlayerStats(player, config.stats_pattern)[category].battles;
  const teamSummary = config.team_summary;

  let minBattles: number;
  switch (category) {
    case "ship":
      minBattles = teamSummary.min_ship_battles;
      break;
    case "overall":
      minBattles = teamSummary.min_overall_battles;
      break;
  }

  return battles >= minBattles;
};

const mean = (players: model.Player[], column: SummaryColumn): number => {
  const values = players.map((player) => column.value(player));
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
