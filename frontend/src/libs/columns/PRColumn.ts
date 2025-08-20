import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { ColorCode } from "@libs/ColorCode";
import { Rating } from "@libs/Rating";
import type { Optional, StatsCategory } from "@libs/types";
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
    const value = this.getPlayerStats(player)[this.category].pr;
    if (value === -1) {
      return undefined;
    }

    return Rating.fromPR(value)?.getColor()?.getFixedTextColor();
  }

  override getDisplayValue(player: data.Player): string {
    const value = this.getPlayerStats(player)[this.category].pr;
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
