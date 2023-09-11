import { DispName } from "src/lib/DispName";
import { AbstractGraphColumn } from "src/lib/column/intetface/AbstractGraphColumn";
import type { OverallKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type {
  StackedBarGraphItem,
  StackedBarGraphParam,
} from "src/lib/value_object/StackedBarGraphParam";
import type { domain } from "wailsjs/go/models";

export class UsingShipTypeRate extends AbstractGraphColumn<OverallKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): OverallKey {
    return "using_ship_type_rate";
  }

  minDisplayName(): string {
    return "艦割合";
  }

  fullDisplayName(): string {
    return "艦種別プレイ割合";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays.overall.using_ship_type_rate;
  }

  displayValue(player: domain.Player): StackedBarGraphParam {
    const digit = this.userConfig.custom_digit.using_ship_type_rate;
    const shipTypeGroup = toPlayerStats(player, this.userConfig.stats_pattern)
      .overall.using_ship_type_rate;
    const ownShipType = player.ship_info.type;
    const colors = this.userConfig.custom_color.ship_type;

    let items: StackedBarGraphItem[] = [];
    DispName.SHIP_TYPES.forEach((st) => {
      const value = shipTypeGroup[st.key];
      const colorCode =
        st.key === ownShipType ? colors.own[st.key] : colors.other[st.key];
      items.push({ label: st.value, colorCode: colorCode, value: value });
    });

    return { digit: digit, items: items };
  }
}
