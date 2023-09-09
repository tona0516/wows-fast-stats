// @ts-ignore
import { ColumnArray } from "src/lib/ColumnArray";
import { AvgTier } from "src/lib/column/AvgTier";
import { Battles } from "src/lib/column/Battles";
import { Damage } from "src/lib/column/Damage";
import { Exp } from "src/lib/column/Exp";
import { HitRate } from "src/lib/column/HitRate";
import { KDRate } from "src/lib/column/KDRate";
import { Kill } from "src/lib/column/Kill";
import { MaxDamage } from "src/lib/column/MaxDamage";
import { PR } from "src/lib/column/PR";
import { PlanesKilled } from "src/lib/column/PlanesKilled";
import { PlayerName } from "src/lib/column/PlayerName";
import { ShipInfo } from "src/lib/column/ShipInfo";
import { SurvivedRate } from "src/lib/column/SurvivedRate";
import { UsingShipTypeRate } from "src/lib/column/UsingShipTypeRate";
import { UsingTierRate } from "src/lib/column/UsingTierRate";
import { WinRate } from "src/lib/column/WinRate";
import type { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import {
  type BasicKey,
  type CommonStatsKey,
  type OptionalBattle,
  type OptionalSummary,
  type OverallKey,
  type OverallOnlyKey,
  type ShipKey,
  type ShipOnlyKey,
  type StatsCategory,
} from "src/lib/types";
import { domain } from "wailsjs/go/models";

export const tableColumns = (
  userConfig: domain.UserConfig,
): [
  basic: ColumnArray<BasicKey>,
  ship: ColumnArray<ShipKey>,
  overall: ColumnArray<OverallKey>,
] => {
  return [
    new ColumnArray<BasicKey>("basic", [
      new PlayerName(userConfig),
      new ShipInfo(userConfig),
    ]),
    new ColumnArray<ShipKey>("ship", [
      new PR(userConfig, "ship"),
      new Damage(userConfig, "ship"),
      new MaxDamage(userConfig, "ship"),
      new WinRate(userConfig, "ship"),
      new KDRate(userConfig, "ship"),
      new Kill(userConfig, "ship"),
      new Exp(userConfig, "ship"),
      new Battles(userConfig, "ship"),
      new SurvivedRate(userConfig, "ship"),
      new PlanesKilled(userConfig),
      new HitRate(userConfig),
    ]),
    new ColumnArray<OverallKey>("overall", [
      new PR(userConfig, "overall"),
      new Damage(userConfig, "overall"),
      new MaxDamage(userConfig, "overall"),
      new WinRate(userConfig, "overall"),
      new KDRate(userConfig, "overall"),
      new Kill(userConfig, "overall"),
      new Exp(userConfig, "overall"),
      new Battles(userConfig, "overall"),
      new SurvivedRate(userConfig, "overall"),
      new AvgTier(userConfig),
      new UsingShipTypeRate(userConfig),
      new UsingTierRate(userConfig),
    ]),
  ];
};

// TODO: 重複コードをマージする
export const displayColumns = (): [
  common: IColumn<CommonStatsKey>[],
  shipOnly: IColumn<ShipOnlyKey>[],
  overallOnly: IColumn<OverallOnlyKey>[],
] => {
  const userConfig = new domain.UserConfig();
  const category: StatsCategory = "ship";

  return [
    [
      new PR(userConfig, category),
      new Damage(userConfig, category),
      new MaxDamage(userConfig, category),
      new WinRate(userConfig, category),
      new KDRate(userConfig, category),
      new Kill(userConfig, category),
      new Exp(userConfig, category),
      new Battles(userConfig, category),
      new SurvivedRate(userConfig, category),
    ],
    [new PlanesKilled(userConfig), new HitRate(userConfig)],
    [
      new AvgTier(userConfig),
      new UsingShipTypeRate(userConfig),
      new UsingTierRate(userConfig),
    ],
  ];
};

export const toPlayerStats = (
  player: domain.Player,
  statsPattern: string,
): domain.PlayerStats => {
  switch (statsPattern) {
    case "pvp_solo":
      return player.pvp_solo;
    case "pvp_all":
      return player.pvp_all;
    default:
      throw Error(`unexpeted error: statsPattern: ${statsPattern}`);
  }
};

export const calculateSummary = (
  battle: OptionalBattle,
  excludes: number[],
  userConfig: domain.UserConfig,
): OptionalSummary => {
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
    overallColspan: columns.filter((it) => it.category() === "overall").length,
    labels: labels,
    friends: friends,
    enemies: enemies,
    diffs: diffs,
  };
};
