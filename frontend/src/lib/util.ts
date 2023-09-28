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
import { type IColumn } from "src/lib/column/intetface/IColumn";
import {
  includesOveralls,
  includesShips,
  type BasicKey,
  type OverallKey,
  type ShipKey,
  type ShipType,
} from "src/lib/types";
import { type DisplayItem } from "src/lib/value_object/DisplayItem";
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

export const displayItems = (): DisplayItem[] => {
  const config = new domain.UserConfig();
  const [_, shipColumns, overallColumns] = tableColumns(config);
  const columns: IColumn<any>[] = [...shipColumns, ...overallColumns];

  let items: DisplayItem[] = [];
  columns.forEach((column) => {
    const displayKey = column.displayKey();
    if (items.filter((it) => it.digitKey === displayKey).length !== 0) {
      return;
    }

    const item: DisplayItem = {
      name: column.fullDisplayName(),
      digitKey: displayKey,
    };

    if (includesShips(displayKey)) {
      item.shipKey = displayKey;
    }

    if (includesOveralls(displayKey)) {
      item.overallKey = displayKey;
    }

    items.push(item);
  });

  return Array.from(items);
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

export const tierString = (value: number): string => {
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
};

export const isShipType = (value: string): value is ShipType => {
  try {
    value as ShipType;
    return true;
  } catch {
    return false;
  }
};
