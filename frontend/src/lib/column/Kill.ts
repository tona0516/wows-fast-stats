import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Kill extends AbstractColumn<CommonKey> implements ISingleColumn {
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super("kill", "撃沈", "平均撃沈数", 1, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].kill;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    const digit = this.config.custom_digit.kill;
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].kill;
    return value.toFixed(digit);
  }

  getTextColorCode(_: domain.Player): string {
    return "";
  }
}
