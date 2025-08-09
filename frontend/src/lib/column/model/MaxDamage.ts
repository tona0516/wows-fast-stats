import MaxDamageTableData from "src/component/stats/internal/table_data/MaxDamageTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { NumbersURL } from "src/lib/NumbersURL";
import type { StatsCategory } from "src/lib/types";
import { tierString } from "src/lib/util";
import type { data } from "wailsjs/go/models";

export interface MaxDamageParam {
  damage: string;
  shipInfo?: {
    url: string;
    name: string;
  };
}

export class MaxDamage extends AbstractStatsColumn<MaxDamageParam> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("max_damage", config, category);
  }

  displayValue(player: data.Player): MaxDamageParam {
    const maxDamage = this.playerStats(player)[this.category].max_damage;
    const value = maxDamage.value.format(this.digit());

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

  getTableDataComponent() {
    return MaxDamageTableData;
  }

  getCssClass(): string | undefined {
    return "text-right";
  }
}
