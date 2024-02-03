import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { type StatsCategory } from "src/lib/types";
import type { model } from "wailsjs/go/models";

export class Kill extends AbstractStatsColumn<string> {
  constructor(config: model.UserConfig, category: StatsCategory) {
    super("kill", 1, config, category);
  }

  displayValue(player: model.Player): string {
    const value = this.playerStats(player)[this.category].kill;
    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
