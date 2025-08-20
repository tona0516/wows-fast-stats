import PlayerNameTableData from "@components/tabledata/PlayerNameTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { Rating } from "@libs/Rating";
import {
  storedPlayerNameColor,
  storedShowClanNation,
  storedStatsExtra,
} from "@libs/stores";
import type { Optional } from "@libs/types";
import type { data } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractColumn } from "./AbstractColumn";

export class PlayerNameColumn extends AbstractColumn {
  constructor() {
    super("player_name", "プレイヤー");
  }

  override needsShow(): boolean {
    return true;
  }

  override getBgColorCode(player: data.Player): Optional<ColorCode> {
    const statsExtra = get(storedStatsExtra);
    const colorPattern = get(storedPlayerNameColor);
    if (colorPattern === "none") {
      return undefined;
    }

    const value = player[statsExtra][colorPattern].pr;
    return Rating.fromPR(value)
      ?.getColor()
      ?.getFixedBgColor(ColorCode.PLAYER_BG_FIDED_RATE);
  }

  override getTableDataComponent() {
    return PlayerNameTableData;
  }

  getClanTag(player: data.Player): string {
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;

    return clanID !== 0 ? `[${clanTag}]` : "";
  }

  getNationFlagClass(player: data.Player): string {
    if (!get(storedShowClanNation)) {
      return "";
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
      return "";
    }

    return `fi fi-${fragIcon}`;
  }

  getPlayerName(player: data.Player): string {
    return player.player_info.name;
  }

  isNPC(player: data.Player): boolean {
    return player.player_info.id === 0;
  }

  getClanColorCode(player: data.Player): Optional<ColorCode> {
    const color = player.player_info.clan.hex_color;

    let colorCode: ColorCode;
    try {
      colorCode = new ColorCode(color);
    } catch {
      return undefined;
    }

    return colorCode.getFixedTextColor();
  }
}
