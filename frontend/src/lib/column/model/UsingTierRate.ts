import StackedBarGraphTableData from "src/component/stats/internal/table_data/StackedBarGraphTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import type { data } from "wailsjs/go/models";

export class UsingTierRate extends AbstractStatsColumn<StackedBarGraphParam> {
  constructor() {
    super("using_tier_rate", "overall");
  }

  displayValue(player: data.Player): StackedBarGraphParam {
    const tierRateGroup = this.playerStats(player).overall.using_tier_rate;
    return {
      digit: this.digit(),
      items: [
        {
          label: "1~4",
          colorCode: "#8CA113",
          value: tierRateGroup.low,
        },
        {
          label: "5~7",
          colorCode: "#205B85",
          value: tierRateGroup.middle,
        },
        {
          label: "8~★",
          colorCode: "#990F4F",
          value: tierRateGroup.high,
        },
      ],
    };
  }

  getTableDataComponent() {
    return StackedBarGraphTableData;
  }
}
