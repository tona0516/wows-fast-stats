import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class AvgTierColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("avgTier", "overall");
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const value = this.getPlayerStats(player).overall.avgTier;
    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
