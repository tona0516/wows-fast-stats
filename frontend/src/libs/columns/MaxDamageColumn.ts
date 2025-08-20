import MaxDamageTableData from "@components/tabledata/MaxDamageTableData.svelte";
import type { StatsCategory } from "@libs/types";
import { shipNumbersURL, toTierString } from "@libs/utils";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export interface MaxDamageParam {
  damage: string;
  shipInfo?: {
    url: string;
    name: string;
  };
}

export class MaxDamageColumn extends AbstractStatsColumn<MaxDamageParam> {
  constructor(category: StatsCategory) {
    super("max_damage", category);
  }

  override getTableDataComponent() {
    return MaxDamageTableData;
  }

  override getDisplayValue(player: data.Player): MaxDamageParam {
    const maxDamage = this.getPlayerStats(player)[this.category].max_damage;
    const value = maxDamage.value.toFixed(this.getDigit());

    switch (this.category) {
      case "ship":
        return { damage: value };
      case "overall": {
        const url = shipNumbersURL(maxDamage.ship_id);
        const name = `${toTierString(maxDamage.ship_tier)} ${maxDamage.ship_name}`;
        return {
          damage: value,
          shipInfo: { url, name },
        };
      }
    }
  }

  override getCssClass(): string {
    return "text-right";
  }
}
