import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { ThreatLevelGenerator } from "src/lib/ThreatLevel";
import type { data } from "wailsjs/go/models";

export class ThreatLevel extends AbstractStatsColumn<string> {
  constructor() {
    super("threat_level", "overall");
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player).overall.threat_level.modified;
    if (value === -1) {
      return "N/A";
    }

    return value.format(this.digit());
  }

  textColorCode(player: data.Player): string {
    return (
      ThreatLevelGenerator.fromScore(player.pvp_all.overall.threat_level.raw)
        ?.color?.text || ""
    );
  }

  getBackgroundColorCode(player: data.Player): string | undefined {
    return ThreatLevelGenerator.fromScore(
      player.pvp_all.overall.threat_level.raw,
    )?.color?.background;
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-right";
  }
}
