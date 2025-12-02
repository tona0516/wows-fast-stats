import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { data } from "@wails/go/models";
import { type ColorCode, THREAT_LEVEL_COLORS } from "../ColorCode";
import type { Optional } from "../types";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class ThreatLevelColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("threat_level", "overall");
  }

  override getTextColorCode(player: data.Player): Optional<ColorCode> {
    return THREAT_LEVEL_COLORS[player.pvp_all.overall.threat_level.rank].text;
  }

  override getBgColorCode(player: data.Player): Optional<ColorCode> {
    return THREAT_LEVEL_COLORS[player.pvp_all.overall.threat_level.rank]
      .background;
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    const value = this.getPlayerStats(player).overall.threat_level.modified;
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
