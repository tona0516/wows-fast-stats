import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import { storedDisplayPref } from "@libs/stores";
import type { StatsCategory } from "@libs/types";
import { formatWithSuffix } from "@libs/utils";
import type { core } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class BattlesColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("battles", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const value = this.getPlayerStats(player)[this.category].battles;
    const pref = get(storedDisplayPref).battles;
    const digit = this.getDigit();

    if (pref.siPrefix) {
      return formatWithSuffix(value, digit);
    }

    return value.toFixed(digit);
  }

  override getCssClass(): string {
    return "text-right";
  }
}
