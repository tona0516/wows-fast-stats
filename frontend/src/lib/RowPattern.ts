import { toPlayerStats } from "src/lib/util";
import type { data } from "wailsjs/go/models";

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
    statsPattern: string,
    allColumnCount: number,
    shipColumnCount: number,
  ): RowPattern => {
    if (allColumnCount === 0) {
      return RowPattern.NO_COLUMN;
    }

    if (player.player_info.is_hidden === true) {
      return RowPattern.PRIVATE;
    }

    const stats = toPlayerStats(player, statsPattern);
    if (player.player_info.id === 0 || stats.overall.battles === 0) {
      return RowPattern.NO_STATS;
    }

    if (stats.ship.battles === 0 && shipColumnCount > 0) {
      return RowPattern.NO_SHIP_STATS;
    }

    return RowPattern.FULL;
  };
}
