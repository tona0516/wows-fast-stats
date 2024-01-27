import PlayerNameTableData from "src/component/main/internal/table_data/PlayerNameTableData.svelte";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { RatingColorFactory } from "src/lib/rating/RatingColorFactory";
import { type BasicKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class PlayerName extends AbstractColumn<BasicKey> {
  constructor(private config: model.UserConfig) {
    super("player_name", "プレイヤー", "プレイヤー", 2);
  }

  getSvelteComponent() {
    return PlayerNameTableData;
  }

  shouldShowColumn(): boolean {
    return true;
  }

  clanName(player: model.Player): string {
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;

    return clanID !== 0 ? `[${clanTag}] ` : "";
  }

  playerName(player: model.Player): string {
    return player.player_info.name;
  }

  isShowCheckBox(player: model.Player): boolean {
    return player.player_info.id !== 0;
  }

  clanColorCode(player: model.Player): string {
    return player.player_info.clan.hex_color
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
      return "#000000";
    }

    const pr = toPlayerStats(player, this.config.stats_pattern)[statsCategory]
      .pr;

    return RatingColorFactory.fromPR(pr, this.config).getTextColorCode();
  }
}
