import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { ColorCode } from "@libs/ColorCode";
import { Rating } from "@libs/Rating";
import type { Optional, StatsCategory } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class WinRateColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("win_rate", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getTextColorCode(player: data.Player): Optional<ColorCode> {
    const value = this.getPlayerStats(player)[this.category].win_rate;
    return Rating.fromWinRate(value)?.getColor()?.getFixedTextColor();
  }

  override getDisplayValue(player: data.Player): string {
    return `${this.getPlayerStats(player)[this.category].win_rate.toFixed(this.getDigit())}%`;
  }

  override getCssClass(): string {
    return "text-right";
  }
}
