import StackedBarGraphTableData from "src/component/stats/internal/table_data/StackedBarGraphTableData.svelte";
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class UsingShipTypeRate extends AbstractStatsColumn<StackedBarGraphParam> {
  constructor(config: data.UserConfigV2) {
    super("using_ship_type_rate", config, "overall");
  }

  displayValue(player: data.Player): StackedBarGraphParam {
    const shipTypeGroup = this.playerStats(player).overall.using_ship_type_rate;
    return {
      digit: this.digit(),
      items: [
        {
          label: "潜水艦",
          colorCode: "#233B8B",
          value: shipTypeGroup.ss,
        },
        {
          label: "駆逐艦",
          colorCode: "#D9760F",
          value: shipTypeGroup.dd,
        },
        {
          label: "巡洋艦",
          colorCode: "#27853F",
          value: shipTypeGroup.cl,
        },
        {
          label: "戦艦",
          colorCode: "#CA1028",
          value: shipTypeGroup.bb,
        },
        {
          label: "空母",
          colorCode: "#5E2883",
          value: shipTypeGroup.cv,
        },
      ],
    };
  }

  svelteComponent() {
    return StackedBarGraphTableData;
  }
}
