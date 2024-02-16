import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { type StatsCategory } from "src/lib/types";
import type { model } from "wailsjs/go/models";

export class KDRate extends AbstractStatsColumn<string> {
  constructor(config: model.UserConfigV2, category: StatsCategory) {
    super("kd_rate", 1, config, category);
  }

  displayValue(player: model.Player): string {
    const value = this.playerStats(player)[this.category].kd_rate;
    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
