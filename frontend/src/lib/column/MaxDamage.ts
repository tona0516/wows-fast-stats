import { BASE_NUMBERS_URL } from "src/const";
import { type IColumn } from "src/lib/column/intetface/IColumn";
import { type CommonStatsKey, type StatsCategory } from "src/lib/types";
import { tierString, toPlayerStats } from "src/lib/util";
import MaxDamageTableData from "src/tabledata_component/MaxDamageTableData.svelte";
import type { domain } from "wailsjs/go/models";

export class MaxDamage implements IColumn<CommonStatsKey> {
  constructor(
    private userConfig: domain.UserConfig,
    private category: StatsCategory,
  ) {}
  displayKey(): CommonStatsKey {
    return "max_damage";
  }

  minDisplayName(): string {
    return "最大Dmg";
  }

  fullDisplayName(): string {
    return "最大ダメージ";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this.category].max_damage;
  }

  countInnerColumn(): number {
    switch (this.category) {
      case "ship":
        return 1;
      case "overall":
        return 2;
    }
  }

  svelteComponent() {
    return MaxDamageTableData;
  }

  damage(player: domain.Player): string {
    const value = toPlayerStats(player, this.userConfig.stats_pattern)[
      this.category
    ].max_damage.damage;
    const digit = this.userConfig.custom_digit.max_damage;
    return value.toFixed(digit);
  }

  shipInfo(player: domain.Player): [url: string, text: string] {
    const maxDamage = toPlayerStats(player, this.userConfig.stats_pattern)
      .overall.max_damage;
    const url =
      BASE_NUMBERS_URL +
      "ship/" +
      maxDamage.ship_id +
      "," +
      maxDamage.ship_name.replaceAll(" ", "-");

    const text = `${tierString(maxDamage.ship_tier)} ${maxDamage.ship_name}`;

    return [url, text];
  }
}
