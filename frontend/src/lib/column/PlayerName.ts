import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { RatingColorFactory } from "src/lib/rating/RatingColorFactory";
import { type BasicKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class PlayerName extends AbstractColumn<BasicKey> {
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
  ) {
    super("player_name", "プレイヤー", "プレイヤー", 2, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return true;
  }

  displayValue(player: domain.Player): string {
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;
    const playerName = player.player_info.name;

    return clanID !== 0 ? `[${clanTag}] ${playerName}` : playerName;
  }

  isShowCheckBox(player: domain.Player): boolean {
    return player.player_info.id !== 0;
  }

  bgColorCode(player: domain.Player): string {
    let statsCategory: StatsCategory | undefined;

    if (this.config.custom_color.player_name === "ship") {
      statsCategory = "ship";
    }
    if (this.config.custom_color.player_name === "overall") {
      statsCategory = "overall";
    }

    if (!statsCategory) {
      return "#000000";
    }

    const pr = toPlayerStats(player, this.config.stats_pattern)[statsCategory]
      .pr;

    return RatingColorFactory.fromPR(pr, this.config).getBgColorCode();
  }
}
