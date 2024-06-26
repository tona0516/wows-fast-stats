import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { type StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class SurvivedRate extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("survived_rate", 1, config, category);
  }

  displayValue(player: data.Player): string {
    const stats = this.playerStats(player)[this.category];
    return `${stats.win_survived_rate.toFixed(
      this.digit(),
    )}% | ${stats.lose_survived_rate.toFixed(this.digit())}%`;
  }

  svelteComponent() {
    return SingleTableData;
  }

  tdClass(_: data.Player): string {
    return CssClass.TD_MULTI;
  }
}
