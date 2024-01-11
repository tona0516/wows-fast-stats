import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import type { ShipKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class HitRate extends AbstractColumn<ShipKey> implements ISingleColumn {
  constructor(private config: model.UserConfig) {
    super("hit_rate", "Hit率(主|魚)", "命中率 (主砲|魚雷)", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays.ship.hit_rate;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_MULTI;
  }

  getDisplayValue(player: model.Player): string {
    const digit = this.config.custom_digit.hit_rate;
    const stats = toPlayerStats(player, this.config.stats_pattern).ship;
    return `${stats.main_battery_hit_rate.toFixed(
      digit,
    )}% | ${stats.torpedoes_hit_rate.toFixed(digit)}%`;
  }

  getTextColorCode(_: model.Player): string {
    return "";
  }
}
