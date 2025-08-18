import MaxDamageTableData from "src/component/stats/internal/table_data/MaxDamageTableData.svelte";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export interface MaxDamageParam {
  damage: string;
  shipInfo?: {
    url: string;
    name: string;
  };
}

export class MaxDamage extends AbstractStatsColumn<MaxDamageParam> {
  constructor(category: StatsCategory) {
    super("max_damage", category);
  }

  displayValue(player: data.Player): MaxDamageParam {
    const maxDamage = this.playerStats(player)[this.category].max_damage;
    const value = maxDamage.value.toFixed(this.digit());

    switch (this.category) {
      case "ship":
        return { damage: value };
      case "overall": {
        const url = AppFunc.shipNumbersURL(maxDamage.ship_id);
        const name = `${AppFunc.toTierString(maxDamage.ship_tier)} ${maxDamage.ship_name}`;
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
