import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class HitRate extends AbstractStatsColumn<string> {
  constructor() {
    super("hit_rate", "ship");
  }

  displayValue(player: data.Player): string {
    const hitRate = this.playerStats(player).ship.hit_rate;

    const main = hitRate.main_battery.format(this.digit());
    const torps = hitRate.torpedoes.format(this.digit());

    return `${main}|${torps}`;
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-center";
  }
}
