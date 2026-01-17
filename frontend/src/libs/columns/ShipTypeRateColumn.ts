import StackedBarGraphTableData from "@components/tabledata/StackedBarGraphTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { SHIP_TYPE_COLORS, SHIP_TYPES } from "@libs/constants";
import type { StackedBarChartParam } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class ShipTypeRateColumn extends AbstractStatsColumn<
  StackedBarChartParam[]
> {
  constructor() {
    super("using_ship_type_rate", "overall");
  }

  override getTableDataComponent() {
    return StackedBarGraphTableData;
  }

  override getDisplayValue(player: core.Player): StackedBarChartParam[] {
    const shipTypeGroup =
      this.getPlayerStats(player).overall.using_ship_type_rate;

    const params: StackedBarChartParam[] = [];
    SHIP_TYPES.forEach((label, shipType) => {
      const color = SHIP_TYPE_COLORS.get(shipType);
      params.push({
        label: label,
        colorCode: color?.getFixedBgColor(ColorCode.SHIP_TYPE_BG_FIXED_RATE),
        value: shipTypeGroup[shipType] || 0,
      });
    });

    return params;
  }
}
