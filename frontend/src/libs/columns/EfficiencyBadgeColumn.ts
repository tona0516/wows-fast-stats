import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { StatsCategory } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

export class EfficiencyBadgeColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("efficiencyBadge", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: core.Player): string {
    const playerStats = this.getPlayerStats(player);

    switch (this.category) {
      case "ship":
        return playerStats.ship.efficiencyBadge;
      case "overall":
        return Object.values(playerStats.overall.efficiencyBadge).join(" | ");
    }
  }

  override getCssClass(): string {
    return "text-center";
  }
}
