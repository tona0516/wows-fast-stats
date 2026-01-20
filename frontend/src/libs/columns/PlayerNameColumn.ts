import PlayerNameTableData from "@components/tabledata/PlayerNameTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { RATING_COLORS, THREAT_LEVEL_COLORS } from "@libs/constants";
import { storedPref } from "@libs/stores";
import type { Optional, StatsExtra } from "@libs/types";
import type { core } from "@wails/go/models";
import { get } from "svelte/store";
import { AbstractColumn } from "./AbstractColumn";

export class PlayerNameColumn extends AbstractColumn {
  constructor() {
    super("player_name", "プレイヤー");
  }

  override needsShow(): boolean {
    return true;
  }

  override getTextColorCode(player: core.Player): Optional<ColorCode> {
    const cfg = get(storedPref);
    const statsExtra = cfg.statsExtra as StatsExtra;
    const pattern = cfg.column.player.colorPattern;

    if (pattern === "none") {
      return undefined;
    }

    switch (pattern) {
      case "pr_ship": {
        const rating = player[statsExtra].ship.pr.rating;
        const code = RATING_COLORS[rating];
        return code ? code.getFixedTextColor() : undefined;
      }
      case "pr_overall": {
        const rating = player[statsExtra].overall.pr.rating;
        const code = RATING_COLORS[rating];
        return code ? code.getFixedTextColor() : undefined;
      }
      case "threatLevel": {
        const threat = player[statsExtra].overall.threatLevel;
        const pair = THREAT_LEVEL_COLORS[threat.rank];
        return pair ? pair.background : undefined;
      }
      default:
        return undefined;
    }
  }

  override getTableDataComponent() {
    return PlayerNameTableData;
  }

  getClanTag(player: core.Player): string {
    const clanID = player.playerInfo.clan.id;
    const clanTag = player.playerInfo.clan.tag;

    return clanID !== 0 ? `[${clanTag}]` : "";
  }

  getNationFlagClass(player: core.Player): string {
    if (!get(storedPref).column.player.enableNationFlag) {
      return "";
    }
    const langMap: Record<string, string> = { ja: "jp", zh: "cn", ko: "kr" };
    const fragIcon = langMap[player.playerInfo.clan.language] ?? "";
    return fragIcon ? `fi fi-${fragIcon}` : "";
  }

  getPlayerName(player: core.Player): string {
    return player.playerInfo.name;
  }

  isNPC(player: core.Player): boolean {
    return player.playerInfo.id === 0;
  }

  getClanColorCode(player: core.Player): Optional<ColorCode> {
    const color = player.playerInfo.clan.colorCode;

    let colorCode: ColorCode;
    try {
      colorCode = new ColorCode(color);
    } catch {
      return undefined;
    }

    return colorCode.getFixedTextColor();
  }
}
