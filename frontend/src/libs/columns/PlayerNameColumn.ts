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
    const statsExtra = cfg.stats_extra as StatsExtra;
    const pattern = cfg.column.player.color_pattern;

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
      case "threat_level": {
        const threat = player[statsExtra].overall.threat_level;
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
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;

    return clanID !== 0 ? `[${clanTag}]` : "";
  }

  getNationFlagClass(player: core.Player): string {
    if (!get(storedPref).column.player.enable_nation_flag) {
      return "";
    }
    const langMap: Record<string, string> = { ja: "jp", zh: "cn", ko: "kr" };
    const fragIcon = langMap[player.player_info.clan.language] ?? "";
    return fragIcon ? `fi fi-${fragIcon}` : "";
  }

  getPlayerName(player: core.Player): string {
    return player.player_info.name;
  }

  isNPC(player: core.Player): boolean {
    return player.player_info.id === 0;
  }

  getClanColorCode(player: core.Player): Optional<ColorCode> {
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
