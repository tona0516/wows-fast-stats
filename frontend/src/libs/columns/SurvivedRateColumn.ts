import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { StatsCategory } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class SurvivedRateColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("survived_rate", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    const sv = this.getPlayerStats(player)[this.category].survived_rate;

    const all = sv.all.toFixed(this.getDigit());
    const win = sv.all.toFixed(this.getDigit());
    const lose = sv.all.toFixed(this.getDigit());

    return `${all}|${win}|${lose}`;
  }

  override getCssClass(): string {
    return "text-center";
  }
}
