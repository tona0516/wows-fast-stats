import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { Rating } from "src/lib/Rating";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { type StatsCategory } from "src/lib/types";
import type { model } from "wailsjs/go/models";

export class PR extends AbstractStatsColumn<string> implements ISummaryColumn {
  constructor(config: model.UserConfigV2, category: StatsCategory) {
    super("pr", 1, config, category);
  }

  displayValue(player: model.Player): string {
    const value = this.value(player);
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }

  tdClass(player: model.Player): string {
    return this.value(player) === -1 ? CssClass.TD_MULTI : CssClass.TD_NUM;
  }

  textColorCode(player: model.Player): string {
    return Rating.fromPR(
      this.value(player),
      this.config.color.skill.text,
    ).colorCode();
  }

  value(player: model.Player): number {
    return this.playerStats(player)[this.category].pr;
  }
}
