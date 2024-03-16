const BASE_URL = "https://asia.wows-numbers.com/";

export namespace NumbersURL {
  export const clan = (clanID: number): string => `${BASE_URL}clan/${clanID},/`;

  export const player = (accountID: number, accountName: string): string =>
    `${BASE_URL}player/${accountID},${accountName}/`;

  export const ship = (shipID: number): string => `${BASE_URL}ship/${shipID},/`;
}
