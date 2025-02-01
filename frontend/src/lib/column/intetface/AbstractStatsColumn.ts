import type { ColumnSetting } from "src/lib/ColumnSetting";
import { CssClass } from "src/lib/CssClass";
import { DispName } from "src/lib/DispName";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { StatsCategory } from "src/lib/types";
import { deriveColumnSetting, toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export abstract class AbstractStatsColumn<T> extends AbstractColumn {
  columnSetting: ColumnSetting;

  constructor(
    readonly key: string,
    readonly innerColumnCount: number,
    readonly config: model.UserConfigV2,
    readonly category: StatsCategory,
  ) {
    super(key, DispName.MIN_COLUMN_NAMES.get(key) ?? key, innerColumnCount);
    this.columnSetting = deriveColumnSetting(config, key);
  }

  abstract displayValue(player: model.Player): T;

  shouldShow(): boolean {
    return this.columnSetting[this.category].value;
  }

  digit(): number {
    return this.columnSetting.digit.value;
  }

  tdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  textColorCode(_: model.Player): string {
    return "";
  }

  bgColorCode(_: model.Player): string {
    return "";
  }

  playerStats(player: model.Player): model.PlayerStats {
    return toPlayerStats(player, this.config.stats_pattern);
  }
}
