import { STATS_COLUMN_INFO } from "@libs/constants";
import { storedOptionalSetting, storedStatsColumnSettings } from "@libs/stores";
import type { StatsCategory, StatsExtra, StatsKey } from "@libs/types";
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
    const cs = get(storedStatsColumnSettings)[this.key];
    switch (this.category) {
      case "ship":
        return cs.is_show_ship;
      case "overall":
        return cs.is_show_overall;
    }
  }

  getDigit(): number {
    return get(storedStatsColumnSettings)[this.key].digit;
  }

  getPlayerStats(player: data.Player): data.PlayerStats {
    return player[get(storedOptionalSetting).stats_extra as StatsExtra];
  }

  getCssClass(): string {
    return "";
  }
}
