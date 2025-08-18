import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
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

    return value.toFixed(this.digit());
  }

  textColorCode(player: data.Player): string {
    const threatLevelInfo = AppFunc.getThreatLevel(
      player.pvp_all.overall.threat_level.raw,
    );
    if (!threatLevelInfo) {
      return "";
    }
    return AppConst.THREAT_LEVEL_COLORS.get(threatLevelInfo.level)?.text || "";
  }

  getBackgroundColorCode(player: data.Player): string | undefined {
    const threatLevelInfo = AppFunc.getThreatLevel(
      player.pvp_all.overall.threat_level.raw,
    );
    if (!threatLevelInfo) {
      return "";
    }
    return (
      AppConst.THREAT_LEVEL_COLORS.get(threatLevelInfo.level)?.background || ""
    );
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-right";
  }
}
