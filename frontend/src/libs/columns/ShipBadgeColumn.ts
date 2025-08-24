import SingleTableData from "@components/tabledata/SingleTableData.svelte";
import type { StatsCategory } from "@libs/types";
import type { data } from "@wails/go/models";
import { AbstractStatsColumn } from "./AbstractStatsColumn";

type ShipBadge = Readonly<data.ShipBadgeGroup>;

const SHIP_BADGES: { [key in keyof ShipBadge]: string } = {
  expert: "E",
  first: "1st",
  second: "2nd",
  third: "3rd",
};

export class ShipBadgeColumn extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("ship_badge", category);
  }

  override getTableDataComponent() {
    return SingleTableData;
  }

  override getDisplayValue(player: data.Player): string {
    const playerStats = this.getPlayerStats(player);

    switch (this.category) {
      case "ship":
        return playerStats.ship.ship_badge;
      case "overall":
        return Object.entries(playerStats.overall.ship_badge)
          .map(
            (entry) =>
              `${SHIP_BADGES[entry[0] as keyof ShipBadge]}:${entry[1]}`,
          )
          .join("|");
    }
  }

  override getCssClass(): string {
    return "text-center";
  }
}
