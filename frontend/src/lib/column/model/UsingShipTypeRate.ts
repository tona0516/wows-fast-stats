import StackedBarGraphTableData from "src/component/main/internal/table_data/StackedBarGraphTableData.svelte";
import { DispName } from "src/lib/DispName";
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class UsingShipTypeRate extends AbstractStatsColumn<StackedBarGraphParam> {
  constructor(config: data.UserConfigV2) {
    super("using_ship_type_rate", 1, config, "overall");
  }

  displayValue(player: data.Player): StackedBarGraphParam {
    const shipTypeGroup = this.playerStats(player).overall.using_ship_type_rate;
    const ownShipType = player.warship.type;
    const colors = this.config.color.ship_type;

    const items = DispName.SHIP_TYPES.toArray().map((it) => {
      const colorCode =
        it.key === ownShipType ? colors.own[it.key] : colors.other[it.key];
      const rate = shipTypeGroup[it.key];

      return { label: it.value, colorCode: colorCode, value: rate };
    });

    return { digit: this.digit(), items: items };
  }

  svelteComponent() {
    return StackedBarGraphTableData;
  }
}
