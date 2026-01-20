import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class HitRateColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("hitRate", "ship");
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const hitRate = this.getPlayerStats(player).ship.hitRate;

    const main = hitRate.mainBattery.toFixed(this.getDigit());
    const torps = hitRate.torpedoes.toFixed(this.getDigit());

    return `${main}% | ${torps}%`;
  }

  override getCssClass(): string {
    return "text-center";
  }
}
