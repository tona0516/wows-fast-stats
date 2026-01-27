import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { ColorCode } from "@libs/ColorCode";
import { RATING_COLORS } from "@libs/constants";
import { storedDisplayPref } from "@libs/stores";
import type { Optional, StatsCategory } from "@libs/types";
import { formatWithSuffix } from "@libs/utils";
import type { core } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class DamageColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("damage", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const value = this.getPlayerStats(player)[this.category].damage.value;
    const pref = get(storedDisplayPref).damage;
    const digit = this.getDigit();

    if (pref.siPrefix) {
      return formatWithSuffix(value, digit);
    }

    return value.toFixed(digit);
  }

  override getTextColorCode(player: core.Player): Optional<ColorCode> {
    if (this.category !== "ship") {
      return undefined;
    }
    const value = this.getPlayerStats(player).ship.damage.rating;

    return RATING_COLORS[value]?.getFixedTextColor();
  }

  override getCssClass(): string {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }
}
