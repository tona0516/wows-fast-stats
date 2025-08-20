import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class PlanesKilledColumn extends AbstractStatsColumn<string> {
  constructor() {
    super("planes_killed", "ship");
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    const value = this.getPlayerStats(player).ship.planes_killed;
    return value.toFixed(this.getDigit());
  }

  override getCssClass(): string {
    return "text-right";
  }
}
