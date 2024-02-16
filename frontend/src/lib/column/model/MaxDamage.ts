import MaxDamageTableData from "src/component/main/internal/table_data/MaxDamageTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { NumbersURL } from "src/lib/NumbersURL";
import { type StatsCategory } from "src/lib/types";
import { tierString } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export interface MaxDamageParam {
  damage: string;
  shipInfo?: {
    url: string;
    name: string;
  };
}

export class MaxDamage extends AbstractStatsColumn<MaxDamageParam> {
  constructor(config: model.UserConfigV2, category: StatsCategory) {
    let innerColumnCount: number;
    switch (category) {
      case "ship":
        innerColumnCount = 1;
        break;
      case "overall":
        innerColumnCount = 2;
        break;
    }

    super("max_damage", innerColumnCount, config, category);
  }

  displayValue(player: model.Player): MaxDamageParam {
    const maxDamage = this.playerStats(player)[this.category].max_damage;
    const value = maxDamage.value.toFixed(this.digit());

    switch (this.category) {
      case "ship":
        return { damage: value };
      case "overall": {
        const url = NumbersURL.ship(maxDamage.ship_id);
        const name = `${tierString(maxDamage.ship_tier)} ${maxDamage.ship_name}`;
        return {
          damage: value,
          shipInfo: { url, name },
        };
      }
    }
  }

  svelteComponent() {
    return MaxDamageTableData;
  }
}
