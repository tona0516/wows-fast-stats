import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { ColorCode } from "@libs/ColorCode";
import { Rating } from "@libs/Rating";
import type { Optional, StatsCategory } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class DamageColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("damage", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    return this.getPlayerStats(player)[this.category].damage.toFixed(
      this.getDigit(),
    );
  }

  override getTextColorCode(player: data.Player): Optional<ColorCode> {
    if (this.category !== "ship") {
      return undefined;
    }
    const value = this.getPlayerStats(player).ship.damage;
    const avgDamage = player.ship_info.avg_damage;

    return Rating.fromShipDamage(value, avgDamage)
      ?.getColor()
      ?.getFixedTextColor();
  }

  override getCssClass(): string {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }
}
