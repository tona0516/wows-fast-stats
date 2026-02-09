import MaxDamageTableData from "@components/tabledata/MaxDamageTableData.svelte";
import { NumbersURL } from "@libs/NumbersURL";
import { storedDisplayPref } from "@libs/stores";
import type { StatsCategory } from "@libs/types";
import { formatWithSuffix, toTierString } from "@libs/utils";
import type { core } from "@wails/go/models";
import { get } from "svelte/store";
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
    super("maxDamage", category);
  }

  override getTableDataComponent() {
    return MaxDamageTableData;
  }

  override getDisplayValue(player: core.Player): MaxDamageParam {
    const maxDamage = this.getPlayerStats(player)[this.category].maxDamage;
    const digit = this.getDigit();

    const pref = get(storedDisplayPref);
    const value = pref.isSiPrefixEnabled
      ? formatWithSuffix(maxDamage.value, digit)
      : maxDamage.value.toFixed(digit);

    switch (this.category) {
      case "ship":
        return { damage: value };
      case "overall": {
        const url = NumbersURL.getShip(maxDamage.shipID);
        const name = `${toTierString(maxDamage.shipTier)} ${maxDamage.shipName}`;
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
