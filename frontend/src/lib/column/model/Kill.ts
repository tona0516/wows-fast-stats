import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class Kill extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("kill", config, category);
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player)[this.category].kill;
    return value.format(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
