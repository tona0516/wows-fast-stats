import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class AvgTierColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("avg_tier", "overall");
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    const value = this.getPlayerStats(player).overall.avg_tier;
    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
