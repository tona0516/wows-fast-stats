import { STATS_COLUMN_INFO } from "@libs/constants";
import type { StatsCategory, StatsExtra, StatsKey } from "@libs/types";
import type { core } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractColumn } from "./AbstractColumn";
import { storedDisplayPref } from "@libs/stores";

export abstract class AbstractStatsColumn<T> extends AbstractColumn {
  constructor(
    readonly key: StatsKey,
    readonly category: StatsCategory,
  ) {
    super(key, STATS_COLUMN_INFO[key].minName ?? key);
  }

  abstract getDisplayValue(player: core.Player): T;

  override needsShow(): boolean {
    const pref = get(storedDisplayPref)[this.key];
    switch (this.category) {
      case "ship":
        return "showShip" in pref ? pref.showShip : false;
      case "overall":
        return "showOverall" in pref ? pref.showOverall : false;
    }
  }

  getDigit(): number {
    const pref = get(storedDisplayPref)[this.key];
    return "digit" in pref ? pref.digit : 0;
  }

  getPlayerStats(player: core.Player): core.PlayerStats {
    return player[get(storedDisplayPref).statsExtra];
  }

  getCssClass(): string {
    return "";
  }
}
