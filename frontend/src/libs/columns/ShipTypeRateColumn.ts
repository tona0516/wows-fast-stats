import StackedBarGraphTableData from "@components/tabledata/StackedBarGraphTableData.svelte";
import { SHIP_TYPE_COLORS } from "@libs/ColorCode";
import { SHIP_TYPES } from "@libs/constants";
import type { StackedBarChartParam } from "@libs/StackedBarChartParam";
import type { data } from "@wails/go/models";
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

  override getDisplayValue(player: data.Player): StackedBarChartParam[] {
    const shipTypeGroup =
      this.getPlayerStats(player).overall.using_ship_type_rate;

    const params: StackedBarChartParam[] = [];
    SHIP_TYPES.forEach((label, shipType) => {
      const color = SHIP_TYPE_COLORS[shipType];
      params.push({
        label: label,
        colorCode: color?.getFixedBgColor(),
        value: shipTypeGroup[shipType] || 0,
      });
    });

    return params;
  }
}
