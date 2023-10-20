import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type ShipKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class PlanesKilled
  extends AbstractColumn<ShipKey>
  implements ISingleColumn
{
  constructor(private config: domain.UserConfig) {
    super("planes_killed", "撃墜", "平均撃墜数", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays.ship.planes_killed;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    const digit = this.config.custom_digit.planes_killed;
    const value = toPlayerStats(player, this.config.stats_pattern).ship
      .planes_killed;
    return value.toFixed(digit);
  }

  getTextColorCode(_: domain.Player): string {
    return "";
  }
}
