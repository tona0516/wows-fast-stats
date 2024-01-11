import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export enum RowPattern {
  NO_COLUMN,
  PRIVATE,
  NO_DATA,
  NO_SHIP_STATS,
  FULL,
}

export namespace RowPattern {
  export const derive = (
    player: model.Player,
    statsPattern: string,
    allColumnNumber: number,
  ): RowPattern => {
    if (allColumnNumber === 0) {
      return RowPattern.NO_COLUMN;
    }

    if (player.player_info.is_hidden === true) {
      return RowPattern.PRIVATE;
    }

    const stats = toPlayerStats(player, statsPattern);
    if (player.player_info.id === 0 || stats.overall.battles === 0) {
      return RowPattern.NO_DATA;
    }

    if (stats.ship.battles === 0) {
      return RowPattern.NO_SHIP_STATS;
    }

    return RowPattern.FULL;
  };
}
