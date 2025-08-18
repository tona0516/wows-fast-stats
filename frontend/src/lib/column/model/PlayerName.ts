import chroma from "chroma-js";
import PlayerNameTableData from "src/component/stats/internal/table_data/PlayerNameTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { RatingInfo } from "src/lib/types";
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

    let ratingInfo: RatingInfo | undefined;
    switch (get(storedPlayerNameColor)) {
      case "ship": {
        ratingInfo = AppFunc.getRatingfromPR(player[statsExtra].ship.pr);
        break;
      }
      case "overall": {
        ratingInfo = AppFunc.getRatingfromPR(player[statsExtra].overall.pr);
        break;
      }
    }

    if (!ratingInfo) {
      return undefined;
    }

    let fixedColor = "";
    const color = AppConst.RATING_COLORS.get(ratingInfo.rating);
    if (color) {
      fixedColor = AppFunc.getFixedColorPair(color).background;
    }

    return fixedColor;
  }
}
