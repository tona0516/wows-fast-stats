import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { ColorCode } from "@libs/ColorCode";
import { RATING_COLORS } from "@libs/constants";
import type { Optional, StatsCategory } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class WinRateColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("winRate", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getTextColorCode(player: core.Player): Optional<ColorCode> {
    const rating = this.getPlayerStats(player)[this.category].winRate.rating;
    return RATING_COLORS[rating]?.getFixedTextColor();
  }

  override getDisplayValue(player: core.Player): string {
    const value = this.getPlayerStats(player)[this.category].winRate.value;
    return `${value.toFixed(this.getDigit())}%`;
  }

  override getCssClass(): string {
    return "text-right";
  }
}
