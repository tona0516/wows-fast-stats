import PlayerNameTableData from "src/component/main/internal/table_data/PlayerNameTableData.svelte";
import { RatingInfo } from "src/lib/RatingLevel";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class PlayerName extends AbstractColumn {
  constructor(private config: model.UserConfigV2) {
    super("player_name", "プレイヤー", 2);
  }

  svelteComponent() {
    return PlayerNameTableData;
  }

  shouldShow(): boolean {
    return true;
  }

  clanTag(player: model.Player): string | undefined {
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;

    return clanID !== 0 ? `[${clanTag}] ` : undefined;
  }

  clanFlagIconClass(player: model.Player): string | undefined {
    if (!this.config.show_language_frag) {
      return undefined;
    }

    let fragIcon = "";
    switch (player.player_info.clan.lang) {
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

  playerName(player: model.Player): string {
    return player.player_info.name;
  }

  isNPC(player: model.Player): boolean {
    return player.player_info.id === 0;
  }

  clanColorCode(player: model.Player): string {
    return player.player_info.clan.hex_color;
  }

  textColorCode(player: model.Player): string {
    let statsCategory: StatsCategory | undefined;

    if (this.config.color.player_name === "ship") {
      statsCategory = "ship";
    }
    if (this.config.color.player_name === "overall") {
      statsCategory = "overall";
    }

    if (!statsCategory) {
      return "";
    }

    const pr = toPlayerStats(player, this.config.stats_pattern)[statsCategory]
      .pr;

    return (
      RatingInfo.fromPR(pr, this.config.color.skill.text)?.textColorCode ?? ""
    );
  }
}
