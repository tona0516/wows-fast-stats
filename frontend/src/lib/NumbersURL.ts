import type { domain } from "wailsjs/go/models";

const BASE_URL = "https://asia.wows-numbers.com/";

export namespace NumbersURL {
  export const clan = (player: domain.Player): string =>
    `${BASE_URL}clan/${player.player_info.clan.id}",${player.player_info.clan.tag}`;

  export const player = (player: domain.Player): string =>
    `${BASE_URL}player/${player.player_info.id}",${player.player_info.name}`;

  export const ship = (id: number, name: string): string =>
    `${BASE_URL}ship/${id}",${name.replaceAll(" ", "-")}`;
}
