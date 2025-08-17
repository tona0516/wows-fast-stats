import { AppConst } from "src/lib/AppConst";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { StatsCategory, StatsKey } from "src/lib/types";
import { storedColumnmSettings, storedStatsExtra } from "src/stores";
import { get } from "svelte/store";
import type { data } from "wailsjs/go/models";

export abstract class AbstractStatsColumn<T> extends AbstractColumn {
  constructor(
    readonly key: StatsKey,
    readonly category: StatsCategory,
  ) {
    super(key, AppConst.STATS_COLUMN_INFO[key].min ?? key);
  }

  abstract displayValue(player: data.Player): T;

  needsShow(): boolean {
    const cs = get(storedColumnmSettings)[this.key];
    switch (this.category) {
      case "ship":
        return cs.ship;
      case "overall":
        return cs.overall;
    }
  }

  digit(): number {
    return get(storedColumnmSettings)[this.key].digit;
  }

  textColorCode(_: data.Player): string {
    return "";
  }

  playerStats(player: data.Player): data.PlayerStats {
    return player[get(storedStatsExtra)];
  }

  getBackgroundColorCode(_player: data.Player): string | undefined {
    return undefined;
  }

  getCssClass(): string | undefined {
    return;
  }
}
