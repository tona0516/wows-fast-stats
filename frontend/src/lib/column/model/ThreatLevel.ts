import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { ThreatLevelInfo } from "src/lib/ThreatLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import type { data } from "wailsjs/go/models";

export class ThreatLevel extends AbstractStatsColumn<string> implements ISummaryColumn {
  constructor(config: data.UserConfigV2) {
    super("threat_level", 1, config, "overall");
  }

  displayValue(player: data.Player): string {
    const value = this.value(player);
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

  value(player: data.Player): number {
    return this.playerStats(player).overall.threat_level.modified;
  }
}
