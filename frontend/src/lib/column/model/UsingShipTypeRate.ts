import { DispName } from "src/lib/DispName";
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { IGraphColumn } from "src/lib/column/intetface/IGraphColumn";
import type { OverallKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class UsingShipTypeRate
  extends AbstractColumn<OverallKey>
  implements IGraphColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
  ) {
    super(
      "using_ship_type_rate",
      "艦割合",
      "艦種別プレイ割合",
      1,
      svelteComponent,
    );
  }

  shouldShowColumn(): boolean {
    return this.config.displays.overall.using_ship_type_rate;
  }

  getGraphParam(player: domain.Player): StackedBarGraphParam {
    const digit = this.config.custom_digit.using_ship_type_rate;
    const shipTypeGroup = toPlayerStats(player, this.config.stats_pattern)
      .overall.using_ship_type_rate;
    const ownShipType = player.ship_info.type;
    const colors = this.config.custom_color.ship_type;

    const items = DispName.SHIP_TYPES.toArray().map((type) => {
      const colorCode =
        type.key === ownShipType
          ? colors.own[type.key]
          : colors.other[type.key];
      const rate = shipTypeGroup[type.key];

      return { label: type.value, colorCode: colorCode, value: rate };
    });

    return { digit: digit, items: items };
  }
}
