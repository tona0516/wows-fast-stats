import StackedBarGraphTableData from "src/component/stats/internal/table_data/StackedBarGraphTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type {
  StackedBarGraphItem,
  StackedBarGraphParam,
} from "src/lib/column/StackedBarGraphParam";
import type { data } from "wailsjs/go/models";

export class UsingTierRate extends AbstractStatsColumn<StackedBarGraphParam> {
  constructor() {
    super("using_tier_rate", "overall");
  }

  displayValue(player: data.Player): StackedBarGraphParam {
    const tierRateGroup = this.playerStats(player).overall.using_tier_rate;

    const items: StackedBarGraphItem[] = [];
    AppConst.TIER_GROUPS.forEach((label, group) => {
      let fixedColor = "";
      const color = AppConst.TIER_GROUP_COLORS.get(group);
      if (color) {
        fixedColor = AppFunc.getFixedColorPair(color).background;
      }

      items.push({
        label: label,
        colorCode: fixedColor,
        value: tierRateGroup[group],
      });
    });

    return {
      digit: this.digit(),
      items: items,
    };
  }

  getTableDataComponent() {
    return StackedBarGraphTableData;
  }
}
