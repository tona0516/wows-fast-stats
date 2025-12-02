import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import { type ColorCode, RATING_COLORS } from "@libs/ColorCode";
import type { Optional, Rating, StatsCategory } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class PRColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("pr", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getTextColorCode(player: data.Player): Optional<ColorCode> {
    const pr = this.getPlayerStats(player)[this.category].pr;
    if (pr.value === -1) {
      return undefined;
    }

    return RATING_COLORS[pr.rating as Rating].getFixedTextColor();
  }

  override getDisplayValue(player: data.Player): string {
    const value = this.getPlayerStats(player)[this.category].pr.value;
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
