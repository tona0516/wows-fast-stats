import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type ShipKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class PlanesKilled
  extends AbstractColumn<ShipKey>
  implements ISingleColumn
{
  constructor(private config: model.UserConfig) {
    super("planes_killed", "撃墜", "平均撃墜数", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.display.ship.planes_killed;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    const digit = this.config.digit.planes_killed;
    const value = toPlayerStats(player, this.config.stats_pattern).ship
      .planes_killed;
    return value.toFixed(digit);
  }

  getTextColorCode(_: model.Player): string {
    return "";
  }
}
