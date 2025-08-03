import chroma from "chroma-js";
import PlayerNameTableData from "src/component/stats/internal/table_data/PlayerNameTableData.svelte";
import { Color } from "src/lib/Color";
import { PlayerNameColor } from "src/lib/DispName";
import { type Rating, RatingfGenerator } from "src/lib/RatingLevel";
import { Theme } from "src/lib/Theme";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { toPlayerStats } from "src/lib/util";
import type { data } from "wailsjs/go/models";

export class PlayerName extends AbstractColumn {
  constructor(private config: data.UserConfigV2) {
    super("player_name", "プレイヤー");
  }

  getTableDataComponent() {
    return PlayerNameTableData;
  }

  needsShow(): boolean {
    return true;
  }

  clanTag(player: data.Player): string | undefined {
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;

    return clanID !== 0 ? `[${clanTag}]` : undefined;
  }

  clanFlagIconClass(player: data.Player): string | undefined {
    if (!this.config.show_language_frag) {
      return undefined;
    }

    let fragIcon = "";
    switch (player.player_info.clan.language) {
      case "ja":
        fragIcon = "jp";
        break;
      case "zh":
        fragIcon = "cn";
        break;
      case "ko":
        fragIcon = "kr";
        break;
    }

    if (!fragIcon) {
      return undefined;
    }

    return `fi fi-${fragIcon}`;
  }

  playerName(player: data.Player): string {
    return player.player_info.name;
  }

  isNPC(player: data.Player): boolean {
    return player.player_info.id === 0;
  }

  clanColorCode(player: data.Player): string {
    const color = player.player_info.clan.hex_color;

    if (Theme.isLighter()) {
      return chroma(color).darken(1.25).hex();
    }

    return color;
  }

  getBackgroundColorCode(player: data.Player): string | undefined {
    const playerStats = toPlayerStats(player, this.config.stats_pattern);

    let rating: Rating | undefined;
    switch (this.config.color.player_name) {
      case PlayerNameColor.SHIP: {
        rating = RatingfGenerator.fromPR(playerStats.ship.pr);
        break;
      }
      case PlayerNameColor.OVERALL: {
        rating = RatingfGenerator.fromPR(playerStats.overall.pr);
        break;
      }
    }

    return rating ? Color.Rating.getFixed(rating.level)?.background : undefined;
  }
}
