import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import { THREAT_LEVEL_COLORS } from "@libs/constants";
import type { core } from "@wails/go/models";
import type { ColorCode } from "../ColorCode";
import type { Optional } from "../types";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class ThreatLevelColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("threat_level", "overall");
  }

  override getTextColorCode(player: core.Player): Optional<ColorCode> {
    return THREAT_LEVEL_COLORS[player.pvp_all.overall.threat_level.rank].text;
  }

  override getBgColorCode(player: core.Player): Optional<ColorCode> {
    return THREAT_LEVEL_COLORS[player.pvp_all.overall.threat_level.rank]
      .background;
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
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
