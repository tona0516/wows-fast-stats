import type { data } from "wailsjs/go/models";
import type { StatsExtra } from "./types";

export enum RowPattern {
  NO_COLUMN = 0,
  PRIVATE = 1,
  NO_STATS = 2,
  NO_SHIP_STATS = 3,
  FULL = 4,
}

export namespace RowPattern {
  export const derive = (
    player: data.Player,
    statsExtra: StatsExtra,
    shipColumnCount: number,
    overallColumnCount: number,
  ): RowPattern => {
    if (shipColumnCount + overallColumnCount === 0) {
      return RowPattern.NO_COLUMN;
    }

    if (player.player_info.is_hidden === true) {
      return RowPattern.PRIVATE;
    }

    const stats = player[statsExtra];
    if (player.player_info.id === 0 || stats.overall.battles === 0) {
      return RowPattern.NO_STATS;
    }

    if (stats.ship.battles === 0 && shipColumnCount > 0) {
      return RowPattern.NO_SHIP_STATS;
    }

    return RowPattern.FULL;
  };

  export const getColumnText = (pattern: RowPattern): string => {
    switch (pattern) {
      case RowPattern.NO_COLUMN:
        return "";
      case RowPattern.PRIVATE:
        return "PRIVAYE";
      case RowPattern.NO_STATS:
        return "N/A";
      case RowPattern.NO_SHIP_STATS:
        return "N/A";
      case RowPattern.FULL:
        return "";
    }
  };
}
