import chroma from "chroma-js";
import PlayerNameTableData from "src/component/stats/internal/table_data/PlayerNameTableData.svelte";
import { AppFunc } from "src/lib/AppFunc";
import { Color } from "src/lib/Color";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { PlayerNameColor } from "src/lib/enums";
import { type Rating, RatingfGenerator } from "src/lib/RatingLevel";
import {
  storedPlayerNameColor,
  storedShowClanNation,
  storedStatsExtra,
} from "src/stores";
import { get } from "svelte/store";
import type { data } from "wailsjs/go/models";

export class PlayerName extends AbstractColumn {
  constructor() {
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
    if (!get(storedShowClanNation)) {
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

    if (AppFunc.isLighter()) {
      return chroma(color).darken(1.25).hex();
    }

    return color;
  }

  getBackgroundColorCode(player: data.Player): string | undefined {
    const statsExtra = get(storedStatsExtra);

    let rating: Rating | undefined;
    switch (get(storedPlayerNameColor)) {
      case PlayerNameColor.SHIP: {
        rating = RatingfGenerator.fromPR(player[statsExtra].ship.pr);
        break;
      }
      case PlayerNameColor.OVERALL: {
        rating = RatingfGenerator.fromPR(player[statsExtra].overall.pr);
        break;
      }
    }

    return rating ? Color.Rating.getFixed(rating.level)?.background : undefined;
  }
}
