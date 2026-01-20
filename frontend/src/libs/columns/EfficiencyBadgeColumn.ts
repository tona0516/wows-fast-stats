import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { StatsCategory } from "@libs/types";
import type { core } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

type EfficiencyBadge = Readonly<core.EfficiencyBadgeGroup>;

const EFFICIENCY_BADGES: { [key in keyof EfficiencyBadge]: string } = {
  expert: "E",
  first: "1st",
  second: "2nd",
  third: "3rd",
};

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
        return Object.entries(playerStats.overall.efficiencyBadge)
          .map(
            (entry) =>
              `${EFFICIENCY_BADGES[entry[0] as keyof EfficiencyBadge]}:${entry[1]}`,
          )
          .join(" | ");
    }
  }

  override getCssClass(): string {
    return "text-center";
  }
}
