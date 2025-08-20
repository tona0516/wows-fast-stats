import StackedBarGraphTableData from "@components/tabledata/StackedBarGraphTableData.svelte";
import { TIER_GROUP_COLORS, TIER_GROUPS } from "@libs/constants";
import type { StackedBarChartParam } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class TierRateColumn extends AbstractStatsColumn<
  StackedBarChartParam[]
> {
  constructor() {
    super("using_tier_rate", "overall");
  }

  override getTableDataComponent() {
    return StackedBarGraphTableData;
  }

  override getDisplayValue(player: data.Player): StackedBarChartParam[] {
    const tierRateGroup = this.getPlayerStats(player).overall.using_tier_rate;

    const params: StackedBarChartParam[] = [];
    TIER_GROUPS.forEach((label, tierGroup) => {
      const color = TIER_GROUP_COLORS.get(tierGroup);

      params.push({
        label: label,
        colorCode: color?.getFixedBgColor(),
        value: tierRateGroup[tierGroup],
      });
    });

    return params;
  }
}
