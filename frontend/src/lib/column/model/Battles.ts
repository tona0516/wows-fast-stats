import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { type StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class Battles extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("battles", 1, config, category);
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player)[this.category].battles;
    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
