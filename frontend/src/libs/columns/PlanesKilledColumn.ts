import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class PlanesKilledColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("planesKilled", "ship");
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const value = this.getPlayerStats(player).ship.planesKilled;
    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
