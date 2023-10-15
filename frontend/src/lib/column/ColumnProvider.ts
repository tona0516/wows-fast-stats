import MaxDamageTableData from "src/component/main/internal/table_data/MaxDamageTableData.svelte";
import PlayerNameTableData from "src/component/main/internal/table_data/PlayerNameTableData.svelte";
import ShipInfoTableData from "src/component/main/internal/table_data/ShipInfoTableData.svelte";
import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import StackedBarGraphTableData from "src/component/main/internal/table_data/StackedBarGraphTableData.svelte";
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
import { domain } from "wailsjs/go/models";

class ColumnArray extends Array<AbstractColumn<any>> {
  constructor(
    private category: ColumnCategory,
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
    config: domain.UserConfig,
  ): [basic: ColumnArray, ship: ColumnArray, overall: ColumnArray] => {
    return [
      new ColumnArray("basic", [
        new PlayerName(PlayerNameTableData, config),
        new ShipInfo(ShipInfoTableData, config),
      ]),
      new ColumnArray("ship", [
        new PR(SingleTableData, config, "ship"),
        new Damage(SingleTableData, config, "ship"),
        new MaxDamage(MaxDamageTableData, config, "ship"),
        new WinRate(SingleTableData, config, "ship"),
        new KDRate(SingleTableData, config, "ship"),
        new Kill(SingleTableData, config, "ship"),
        new Exp(SingleTableData, config, "ship"),
        new Battles(SingleTableData, config, "ship"),
        new SurvivedRate(SingleTableData, config, "ship"),
        new PlanesKilled(SingleTableData, config),
        new HitRate(SingleTableData, config),
      ]),
      new ColumnArray("overall", [
        new PR(SingleTableData, config, "overall"),
        new Damage(SingleTableData, config, "overall"),
        new MaxDamage(MaxDamageTableData, config, "overall"),
        new WinRate(SingleTableData, config, "overall"),
        new KDRate(SingleTableData, config, "overall"),
        new Kill(SingleTableData, config, "overall"),
        new Exp(SingleTableData, config, "overall"),
        new Battles(SingleTableData, config, "overall"),
        new SurvivedRate(SingleTableData, config, "overall"),
        new AvgTier(SingleTableData, config),
        new UsingShipTypeRate(StackedBarGraphTableData, config),
        new UsingTierRate(StackedBarGraphTableData, config),
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
    const config = new domain.UserConfig();
    const [_, shipColumns, overallColumns] = getAllColumns(config);
    const columns: AbstractColumn<any>[] = [...shipColumns, ...overallColumns];

    let result: DisplayableColumn[] = [];
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
