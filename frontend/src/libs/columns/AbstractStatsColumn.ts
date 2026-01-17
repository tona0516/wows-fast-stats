import { STATS_COLUMN_INFO } from "@libs/constants";
import { storedPref } from "@libs/stores";
import type { StatsCategory, StatsExtra, StatsKey } from "@libs/types";
import type { core } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractColumn } from "./AbstractColumn";

export abstract class AbstractStatsColumn<T> extends AbstractColumn {
  constructor(
    readonly key: StatsKey,
    readonly category: StatsCategory,
  ) {
    super(key, STATS_COLUMN_INFO[key].min ?? key);
  }

  abstract getDisplayValue(player: core.Player): T;

  override needsShow(): boolean {
    const cs = get(storedPref).column.stats[this.key];
    switch (this.category) {
      case "ship":
        return cs.is_show_ship;
      case "overall":
        return cs.is_show_overall;
    }
  }

  getDigit(): number {
    return get(storedPref).column.stats[this.key].digit;
  }

  getPlayerStats(player: core.Player): core.PlayerStats {
    return player[get(storedPref).stats_extra as StatsExtra];
  }

  getCssClass(): string {
    return "";
  }
}
