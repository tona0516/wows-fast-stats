import MaxDamageTableData from "src/component/main/internal/table_data/MaxDamageTableData.svelte";
import PlayerNameTableData from "src/component/main/internal/table_data/PlayerNameTableData.svelte";
import ShipInfoTableData from "src/component/main/internal/table_data/ShipInfoTableData.svelte";
import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import StackedBarGraphTableData from "src/component/main/internal/table_data/StackedBarGraphTableData.svelte";
import { AvgTier } from "src/lib/column/AvgTier";
import { Battles } from "src/lib/column/Battles";
import { ColumnArray } from "src/lib/column/ColumnArray";
import { Damage } from "src/lib/column/Damage";
import { Exp } from "src/lib/column/Exp";
import { HitRate } from "src/lib/column/HitRate";
import type { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { KDRate } from "src/lib/column/KDRate";
import { Kill } from "src/lib/column/Kill";
import { MaxDamage } from "src/lib/column/MaxDamage";
import { PlanesKilled } from "src/lib/column/PlanesKilled";
import { PlayerName } from "src/lib/column/PlayerName";
import { PR } from "src/lib/column/PR";
import { ShipInfo } from "src/lib/column/ShipInfo";
import { SurvivedRate } from "src/lib/column/SurvivedRate";
import { UsingShipTypeRate } from "src/lib/column/UsingShipTypeRate";
import { UsingTierRate } from "src/lib/column/UsingTierRate";
import { WinRate } from "src/lib/column/WinRate";
import {
  isOverallKey,
  isShipKey,
  type DigitKey,
  type OverallKey,
  type ShipKey,
} from "src/lib/types";
import { domain } from "wailsjs/go/models";

interface DisplayItem {
  name: string;
  digitKey: DigitKey;
  shipKey?: ShipKey;
  overallKey?: OverallKey;
}

export namespace ColumnProvider {
  export const tableColumns = (
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

  export const displayItems = (): DisplayItem[] => {
    const config = new domain.UserConfig();
    const [_, shipColumns, overallColumns] = tableColumns(config);
    const columns: AbstractColumn<any>[] = [...shipColumns, ...overallColumns];

    let items: DisplayItem[] = [];
    columns.forEach((column) => {
      const displayKey = column.getDisplayKey();
      if (items.filter((it) => it.digitKey === displayKey).length !== 0) {
        return;
      }

      const item: DisplayItem = {
        name: column.getFullDisplayName(),
        digitKey: displayKey,
      };

      if (isShipKey(displayKey)) {
        item.shipKey = displayKey;
      }

      if (isOverallKey(displayKey)) {
        item.overallKey = displayKey;
      }

      items.push(item);
    });

    return Array.from(items);
  };
}
