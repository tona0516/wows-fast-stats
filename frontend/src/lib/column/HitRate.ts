import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { CssClass, type ShipOnlyKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class HitRate extends AbstractSingleColumn<ShipOnlyKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): ShipOnlyKey {
    return "hit_rate";
  }

  minDisplayName(): string {
    return "Hit率(主|魚)";
  }

  fullDisplayName(): string {
    return "命中率 (主砲|魚雷)";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays.ship.hit_rate;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_MULTI;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.hit_rate;
    const stats = toPlayerStats(player, this.userConfig.stats_pattern).ship;
    return `${stats.main_battery_hit_rate.toFixed(
      digit,
    )}% | ${stats.torpedoes_hit_rate.toFixed(digit)}%`;
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
