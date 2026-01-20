import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { StatsCategory } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class SurvivedRateColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("survivedRate", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const sv = this.getPlayerStats(player)[this.category].survivedRate;

    const all = sv.all.toFixed(this.getDigit());
    const win = sv.win.toFixed(this.getDigit());
    const lose = sv.lose.toFixed(this.getDigit());

    return `${all}% | ${win}% | ${lose}%`;
  }

  override getCssClass(): string {
    return "text-center";
  }
}
