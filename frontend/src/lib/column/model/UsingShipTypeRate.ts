import StackedBarGraphTableData from "src/component/stats/internal/table_data/StackedBarGraphTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type {
  StackedBarGraphItem,
  StackedBarGraphParam,
} from "src/lib/column/StackedBarGraphParam";
import type { data } from "wailsjs/go/models";

export class UsingShipTypeRate extends AbstractStatsColumn<StackedBarGraphParam> {
  constructor() {
    super("using_ship_type_rate", "overall");
  }

  displayValue(player: data.Player): StackedBarGraphParam {
    const shipTypeGroup = this.playerStats(player).overall.using_ship_type_rate;

    const items: StackedBarGraphItem[] = [];
    for (const type of AppConst.SHIP_TYPES.keys()) {
      let fixedColor = "";
      const color = AppConst.SHIP_TYPE_COLORS.get(type);
      if (color) {
        fixedColor = AppFunc.getFixedColorPair(color).background;
      }

      items.push({
        label: AppConst.SHIP_TYPES.get(type) || "",
        colorCode: fixedColor,
        value: shipTypeGroup[type] || 0,
      });
    }

    return {
      digit: this.digit(),
      items: items,
    };
  }

  getTableDataComponent() {
    return StackedBarGraphTableData;
  }
}
