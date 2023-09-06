import { DispName } from "src/lib/DispName";
import type {
  StackedBarGraphItem,
  StackedBarGraphParam,
} from "src/lib/StackedBarGraphParam";
import { AbstractGraphColumn } from "src/lib/column/intetface/AbstractGraphColumn";
import type { OverallOnlyKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class UsingShipTypeRate extends AbstractGraphColumn<OverallOnlyKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): OverallOnlyKey {
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
    DispName.SHIP_TYPES.forEach((pair) => {
      const value = shipTypeGroup[pair.first];
      const colorCode =
        pair.first === ownShipType
          ? colors.own[pair.first]
          : colors.other[pair.first];
      items.push({ label: pair.second, colorCode: colorCode, value: value });
    });

    return { digit: digit, items: items };
  }
}
