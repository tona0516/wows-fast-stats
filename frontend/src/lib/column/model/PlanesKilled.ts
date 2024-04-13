import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class PlanesKilled extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2) {
    super("planes_killed", 1, config, "ship");
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player).ship.planes_killed;
    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
