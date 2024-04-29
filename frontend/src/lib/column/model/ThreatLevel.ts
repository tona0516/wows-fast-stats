import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { ThreatLevelInfo } from "src/lib/ThreatLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class ThreatLevel extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2) {
    super("threat_level", 1, config, "overall");
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player).overall.threat_level.modified;
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.digit());
  }

  textColorCode(player: data.Player): string {
    return ThreatLevelInfo.fromScore(player.pvp_all.overall.threat_level.raw)
      ?.textColorCode ?? "";
  }

  bgColorCode(player: data.Player): string {
    return ThreatLevelInfo.fromScore(player.pvp_all.overall.threat_level.raw)
      ?.bgColorCode ?? "";
  }

  svelteComponent() {
    return SingleTableData;
  }
}
