import MaxDamageTableData from "src/component/main/internal/table_data/MaxDamageTableData.svelte";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { NumbersURL } from "src/lib/NumbersURL";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { tierString, toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class MaxDamage extends AbstractColumn<CommonKey> {
  constructor(
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    let innerColumnNumber: number;
    switch (category) {
      case "ship":
        innerColumnNumber = 1;
        break;
      case "overall":
        innerColumnNumber = 2;
        break;
    }

    super("max_damage", "最大Dmg", "最大ダメージ", innerColumnNumber);
  }

  getSvelteComponent() {
    return MaxDamageTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].max_damage;
  }

  damage(player: domain.Player): string {
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].max_damage.damage;
    const digit = this.config.custom_digit.max_damage;
    return value.toFixed(digit);
  }

  shipInfo(player: domain.Player): [url: string, text: string] {
    const maxDamage = toPlayerStats(player, this.config.stats_pattern).overall
      .max_damage;
    const url = NumbersURL.ship(maxDamage.ship_id, maxDamage.ship_name);

    const text = `${tierString(maxDamage.ship_tier)} ${maxDamage.ship_name}`;

    return [url, text];
  }
}
