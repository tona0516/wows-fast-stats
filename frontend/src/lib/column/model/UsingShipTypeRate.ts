import StackedBarGraphTableData from "src/component/stats/internal/table_data/StackedBarGraphTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { Color } from "src/lib/Color";
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
      items.push({
        label: AppConst.SHIP_TYPES.get(type) || "",
        colorCode: Color.ShipType.getDefault(type) || "",
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
