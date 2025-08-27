import PlayerNameTableData from "@components/tabledata/PlayerNameTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { Rating } from "@libs/Rating";
import {
  storedPlayerNameColumnSettings,
  storedShipInfoColumnSettings,
  storedStatsExtra,
} from "@libs/stores";
import { ThreatLevel } from "@libs/ThreatLevel";
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

  override getTextColorCode(player: data.Player): Optional<ColorCode> {
    const statsExtra = get(storedStatsExtra);
    const colorPattern = get(storedPlayerNameColumnSettings).colorPattern;

    switch (colorPattern) {
      case "pr_ship": {
        const value = player[statsExtra].ship.pr;
        const rating = Rating.fromPR(value);
        return rating?.getColor().getFixedTextColor();
      }
      case "pr_overall": {
        const value = player[statsExtra].overall.pr;
        const rating = Rating.fromPR(value);
        return rating?.getColor().getFixedTextColor();
      }
      case "threat_level": {
        const value = player[statsExtra].overall.threat_level;
        const level = ThreatLevel.fromScore(value.raw);
        return level?.getColor().text;
      }
      case "none":
        return undefined;
      default:
        return undefined;
    }
  }

  override getBgColorCode(player: data.Player): Optional<ColorCode> {
    const statsExtra = get(storedStatsExtra);
    const colorPattern = get(storedPlayerNameColumnSettings).colorPattern;

    switch (colorPattern) {
      case "pr_ship": {
        const value = player[statsExtra].ship.pr;
        const rating = Rating.fromPR(value);
        return rating?.getColor().getFixedBgColor();
      }
      case "pr_overall": {
        const value = player[statsExtra].overall.pr;
        const rating = Rating.fromPR(value);
        return rating?.getColor().getFixedBgColor();
      }
      case "threat_level": {
        const value = player[statsExtra].overall.threat_level;
        const level = ThreatLevel.fromScore(value.raw);
        return level?.getColor().background;
      }
      case "none":
        return undefined;
      default:
        return undefined;
    }
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
    if (!get(storedShipInfoColumnSettings).enableNationFlag) {
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
