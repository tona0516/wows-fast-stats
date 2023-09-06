import { BASE_NUMBERS_URL } from "src/const";
import type { IColumn } from "src/lib/column/intetface/IColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import type { BasicKey, StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import PlayerNameTableData from "src/tabledata_component/PlayerNameTableData.svelte";
import type { domain } from "wailsjs/go/models";

export class PlayerName implements IColumn<BasicKey> {
  constructor(private userConfig: domain.UserConfig) {}

  displayKey(): BasicKey {
    return "player_name";
  }

  minDisplayName(): string {
    return "プレイヤー";
  }

  fullDisplayName(): string {
    return "プレイヤー";
  }

  shouldShowColumn(): boolean {
    return true;
  }

  countInnerColumn(): number {
    return 1;
  }

  svelteComponent() {
    return PlayerNameTableData;
  }

  displayValue(player: domain.Player): string {
    const clanID = player.player_info.clan.id;
    const clanTag = player.player_info.clan.tag;
    const playerName = player.player_info.name;

    return clanID !== 0 ? `[${clanTag}] ${playerName}` : playerName;
  }

  bgColorCode(player: domain.Player): string {
    let statsCategory: StatsCategory | undefined;

    if (this.userConfig.custom_color.player_name === "ship") {
      statsCategory = "ship";
    }
    if (this.userConfig.custom_color.player_name === "overall") {
      statsCategory = "overall";
    }

    if (!statsCategory) {
      return "";
    }

    const pr = toPlayerStats(player, this.userConfig.stats_pattern)[
      statsCategory
    ].pr;

    return RatingConverterFactory.fromPR(pr, this.userConfig).bgColorCode();
  }

  clanURL(player: domain.Player): string {
    return (
      BASE_NUMBERS_URL +
      "clan/" +
      player.player_info.clan.id +
      "," +
      player.player_info.clan.tag
    );
  }

  playerURL(player: domain.Player): string {
    return (
      BASE_NUMBERS_URL +
      "player/" +
      player.player_info.id +
      "," +
      player.player_info.name
    );
  }
}
