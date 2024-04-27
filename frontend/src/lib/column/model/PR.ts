import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { RatingInfo } from "src/lib/RatingLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { type StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class PR extends AbstractStatsColumn<string> implements ISummaryColumn {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("pr", 1, config, category);
  }

  displayValue(player: data.Player): string {
    const value = this.value(player);
    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }

  tdClass(player: data.Player): string {
    return this.value(player) === -1 ? CssClass.TD_MULTI : CssClass.TD_NUM;
  }

  textColorCode(player: data.Player): string {
    return (
      RatingInfo.fromPR(this.value(player), this.config.color.skill.text)
        ?.textColorCode ?? ""
    );
  }

  value(player: data.Player): number {
    return this.playerStats(player)[this.category].pr;
  }
}
