import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { StatsCategory } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class KillColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("kill", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const value = this.getPlayerStats(player)[this.category].kill;
    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
