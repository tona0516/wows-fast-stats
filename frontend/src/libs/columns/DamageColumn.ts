import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { ColorCode } from "@libs/ColorCode";
import { RATING_COLORS } from "@libs/constants";
import type { Optional, StatsCategory } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class DamageColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("damage", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    return this.getPlayerStats(player)[this.category].damage.value.toFixed(
      this.getDigit(),
    );
  }

  override getTextColorCode(player: data.Player): Optional<ColorCode> {
    if (this.category !== "ship") {
      return undefined;
    }
    const value = this.getPlayerStats(player).ship.damage.rating;

    return RATING_COLORS[value].getFixedTextColor();
  }

  override getCssClass(): string {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }
}
