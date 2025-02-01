import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { model } from "wailsjs/go/models";

export class HitRate extends AbstractStatsColumn<string> {
  constructor(config: model.UserConfigV2) {
    super("hit_rate", 1, config, "ship");
  }

  displayValue(player: model.Player): string {
    const stats = this.playerStats(player).ship;
    return `${stats.main_battery_hit_rate.toFixed(
      this.digit(),
    )}% | ${stats.torpedoes_hit_rate.toFixed(this.digit())}%`;
  }

  svelteComponent() {
    return SingleTableData;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_MULTI;
  }
}
