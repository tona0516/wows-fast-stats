import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { CssClass, type ShipOnlyKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class PlanesKilled extends AbstractSingleColumn<ShipOnlyKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): ShipOnlyKey {
    return "planes_killed";
  }

  minDisplayName(): string {
    return "撃墜";
  }

  fullDisplayName(): string {
    return "平均撃墜数";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays.ship.planes_killed;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.battles;
    const value = toPlayerStats(player, this.userConfig.stats_pattern).ship
      .planes_killed;
    return value.toFixed(digit);
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
