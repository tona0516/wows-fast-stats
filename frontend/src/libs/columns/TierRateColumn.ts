import StackedBarGraphTableData from "@components/tabledata/StackedBarGraphTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import type { StackedBarChartParam } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

type TierGroup = Readonly<keyof core.TierGroup>;

const DISPLAY_NAMES: { [tierGroup in TierGroup]: string } = {
  low: "1~4",
  middle: "5~7",
  high: "8~★",
} as const;

const COLORS: { [tierGroup in TierGroup]: ColorCode } = {
  low: new ColorCode("#8CA113"),
  middle: new ColorCode("#205B85"),
  high: new ColorCode("#990F4F"),
} as const;

export class TierRateColumn extends AbstractStatsColumn<
  StackedBarChartParam[]
> {
  constructor() {
    super("using_tier_rate", "overall");
  }

  override getTableDataComponent() {
    return StackedBarGraphTableData;
  }

  override getDisplayValue(player: core.Player): StackedBarChartParam[] {
    const tierRateGroup = this.getPlayerStats(player).overall.using_tier_rate;

    const params: StackedBarChartParam[] = [];
    Object.keys(DISPLAY_NAMES).forEach((key) => {
      const tierGroup = key as TierGroup;

      const label = DISPLAY_NAMES[tierGroup];
      const color = COLORS[tierGroup];
      params.push({
        label: label,
        colorCode: color.getFixedBgColor(),
        value: tierRateGroup[tierGroup],
      });
    });
    return params;
  }
}
