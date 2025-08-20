import { STATS_COLUMN_INFO } from "@libs/constants";
import { storedColumnmSettings, storedStatsExtra } from "@libs/stores";
import type { StatsCategory, StatsKey } from "@libs/types";
import type { data } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractColumn } from "./AbstractColumn";

export abstract class AbstractStatsColumn<T> extends AbstractColumn {
  constructor(
    readonly key: StatsKey,
    readonly category: StatsCategory,
  ) {
    super(key, STATS_COLUMN_INFO[key].min ?? key);
  }

  abstract getDisplayValue(player: data.Player): T;

  override needsShow(): boolean {
    const cs = get(storedColumnmSettings)[this.key];
    switch (this.category) {
      case "ship":
        return cs.ship;
      case "overall":
        return cs.overall;
    }
  }

  getDigit(): number {
    return get(storedColumnmSettings)[this.key].digit;
  }

  getPlayerStats(player: data.Player): data.PlayerStats {
    return player[get(storedStatsExtra)];
  }

  getCssClass(): string {
    return "";
  }
}
