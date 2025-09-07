import PlayerNameTableData from "@components/tabledata/PlayerNameTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { RATING_COLORS } from "@libs/constants";
import { storedBasicColumnSetting, storedOptionalSetting } from "@libs/stores";
import { ThreatLevel } from "@libs/ThreatLevel";
import type { Optional, StatsExtra } from "@libs/types";
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
    const statsExtra = get(storedOptionalSetting).stats_extra as StatsExtra;
    const colorPattern = get(storedBasicColumnSetting).player.color_pattern;

    switch (colorPattern) {
      case "pr_ship": {
        const value = player[statsExtra].ship.pr.rating;
        return RATING_COLORS[value].getFixedTextColor();
      }
      case "pr_overall": {
        const value = player[statsExtra].overall.pr.rating;
        return RATING_COLORS[value].getFixedTextColor();
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
    const statsExtra = get(storedOptionalSetting).stats_extra as StatsExtra;
    const colorPattern = get(storedBasicColumnSetting).player.color_pattern;

    switch (colorPattern) {
      case "pr_ship": {
        const value = player[statsExtra].ship.pr.rating;
        return RATING_COLORS[value].getFixedBgColor();
      }
      case "pr_overall": {
        const value = player[statsExtra].overall.pr.rating;
        return RATING_COLORS[value].getFixedBgColor();
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
    if (!get(storedBasicColumnSetting).player.enable_nation_flag) {
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
