import PlayerNameTableData from "src/component/stats/internal/table_data/PlayerNameTableData.svelte";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { data } from "wailsjs/go/models";

export class PlayerName extends AbstractColumn {
  constructor(private config: data.UserConfigV2) {
    super("player_name", "プレイヤー");
  }

  svelteComponent() {
    return PlayerNameTableData;
  }

  shouldShow(): boolean {
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
    return player.player_info.clan.hex_color;
  }

  getBackgroundColorCode(player: data.Player): string | undefined {
    return undefined;
  }
}
