import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class PlanesKilled extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2) {
    super("planes_killed", config, "ship");
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player).ship.planes_killed;
    return value.format(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
