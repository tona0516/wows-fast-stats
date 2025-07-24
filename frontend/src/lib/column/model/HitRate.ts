import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class HitRate extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2) {
    super("hit_rate", config, "ship");
  }

  displayValue(player: data.Player): string {
    const stats = this.playerStats(player).ship;
    const main = stats.main_battery_hit_rate.format(this.digit());
    const torps = stats.torpedoes_hit_rate.format(this.digit());

    return `${main}% | ${torps}%`;
  }

  svelteComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-center";
  }
}
