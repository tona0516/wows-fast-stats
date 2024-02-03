import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { Rating } from "src/lib/Rating";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { type StatsCategory } from "src/lib/types";
import type { model } from "wailsjs/go/models";

export class Damage
  extends AbstractStatsColumn<string>
  implements ISummaryColumn
{
  constructor(config: model.UserConfig, category: StatsCategory) {
    super("damage", 1, config, category);
  }

  displayValue(player: model.Player): string {
    return this.value(player).toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }

  textColorCode(player: model.Player): string {
    if (this.category !== "ship") return "";
    const value = this.playerStats(player).ship.damage;

    return Rating.fromDamage(
      value,
      player.ship_info.avg_damage,
      this.config.color.skill.text,
    ).colorCode();
  }

  value(player: model.Player): number {
    return this.playerStats(player)[this.category].damage;
  }
}
