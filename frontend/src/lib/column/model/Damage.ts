import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { RatingInfo } from "src/lib/RatingLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { type StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class Damage
  extends AbstractStatsColumn<string>
  implements ISummaryColumn
{
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("damage", 1, config, category);
  }

  displayValue(player: data.Player): string {
    return this.value(player).toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    if (this.category !== "ship") return "";
    const value = this.playerStats(player).ship.damage;

    return (
      RatingInfo.fromDamage(
        value,
        player.warship.average_damage,
        this.config.color.skill.text,
      )?.textColorCode ?? ""
    );
  }

  value(player: data.Player): number {
    return this.playerStats(player)[this.category].damage;
  }
}
