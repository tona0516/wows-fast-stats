import StackedBarGraphTableData from "@components/tabledata/StackedBarGraphTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { SHIP_TYPE_COLORS } from "@libs/constants";
import type { ShipType, StackedBarChartParam } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

const SHIP_TYPES: Readonly<Map<ShipType, string>> = new Map<ShipType, string>([
  ["ss", "潜水艦"],
  ["dd", "駆逐艦"],
  ["cl", "巡洋艦"],
  ["bb", "戦艦"],
  ["cv", "空母"],
]);

export class ShipTypeRateColumn extends AbstractStatsColumn<
  StackedBarChartParam[]
> {
  constructor() {
    super("usingShipTypeRate", "overall");
  }

  override getTableDataComponent() {
    return StackedBarGraphTableData;
  }

  override getDisplayValue(player: core.Player): StackedBarChartParam[] {
    const shipTypeGroup = this.getPlayerStats(player).overall.usingShipTypeRate;

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
