import { DispName } from "src/lib/DispName";
import type { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { AvgTier } from "src/lib/column/model/AvgTier";
import { Battles } from "src/lib/column/model/Battles";
import { Damage } from "src/lib/column/model/Damage";
import { Exp } from "src/lib/column/model/Exp";
import { HitRate } from "src/lib/column/model/HitRate";
import { KDRate } from "src/lib/column/model/KDRate";
import { Kill } from "src/lib/column/model/Kill";
import { MaxDamage } from "src/lib/column/model/MaxDamage";
import { PR } from "src/lib/column/model/PR";
import { PlanesKilled } from "src/lib/column/model/PlanesKilled";
import { PlatoonRate } from "src/lib/column/model/PlatoonRate";
import { PlayerName } from "src/lib/column/model/PlayerName";
import { ShipInfo } from "src/lib/column/model/ShipInfo";
import { SurvivedRate } from "src/lib/column/model/SurvivedRate";
import { UsingShipTypeRate } from "src/lib/column/model/UsingShipTypeRate";
import { UsingTierRate } from "src/lib/column/model/UsingTierRate";
import { WinRate } from "src/lib/column/model/WinRate";
import {
  type ColumnCategory,
  type DigitKey,
  type OverallKey,
  type ShipKey,
} from "src/lib/types";
import { isDigitKey, isOverallKey, isShipKey } from "src/lib/util";
import { model } from "wailsjs/go/models";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
class ColumnArray extends Array<AbstractColumn<any>> {
  constructor(
    private category: ColumnCategory,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    private columns: AbstractColumn<any>[],
  ) {
    super(...columns);
  }

  dispName(): string {
    return DispName.COLUMN_CATEGORIES.get(this.category)!;
  }

  columnCount(): number {
    return this.columns
      .filter((it) => it.shouldShowColumn())
      .reduce((a, it) => a + it.innerColumnNumber, 0);
  }
}

export namespace ColumnProvider {
  export const getAllColumns = (
    config: model.UserConfig,
  ): [basic: ColumnArray, ship: ColumnArray, overall: ColumnArray] => {
    return [
      new ColumnArray("basic", [new PlayerName(config), new ShipInfo(config)]),
      new ColumnArray("ship", [
        new PR(config, "ship"),
        new Damage(config, "ship"),
        new MaxDamage(config, "ship"),
        new WinRate(config, "ship"),
        new KDRate(config, "ship"),
        new Kill(config, "ship"),
        new Exp(config, "ship"),
        new Battles(config, "ship"),
        new SurvivedRate(config, "ship"),
        new PlatoonRate(config, "ship"),
        new PlanesKilled(config),
        new HitRate(config),
      ]),
      new ColumnArray("overall", [
        new PR(config, "overall"),
        new Damage(config, "overall"),
        new MaxDamage(config, "overall"),
        new WinRate(config, "overall"),
        new KDRate(config, "overall"),
        new Kill(config, "overall"),
        new Exp(config, "overall"),
        new Battles(config, "overall"),
        new SurvivedRate(config, "overall"),
        new PlatoonRate(config, "overall"),
        new AvgTier(config),
        new UsingShipTypeRate(config),
        new UsingTierRate(config),
      ]),
    ];
  };

  interface DisplayableColumn {
    name: string;
    digitKey?: DigitKey;
    shipKey?: ShipKey;
    overallKey?: OverallKey;
  }

  export const getDisplayableColumns = (): DisplayableColumn[] => {
    const config = new model.UserConfig();
    const [_, shipColumns, overallColumns] = getAllColumns(config);
    const columns: AbstractColumn<string>[] = [
      ...shipColumns,
      ...overallColumns,
    ];

    const result: DisplayableColumn[] = [];
    columns.forEach((column) => {
      const key = column.displayKey;
      if (result.find((it) => it.digitKey === key)) {
        return;
      }

      const dc: DisplayableColumn = { name: column.fullDisplayName };

      if (isDigitKey(key)) dc.digitKey = key;
      if (isShipKey(key)) dc.shipKey = key;
      if (isOverallKey(key)) dc.overallKey = key;

      result.push(dc);
    });

    return result;
  };
}
