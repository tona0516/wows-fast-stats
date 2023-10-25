const BASE_URL = "https://asia.wows-numbers.com/";

export namespace NumbersURL {
  export const clan = (clanID: number): string => `${BASE_URL}clan/${clanID},/`;

  export const player = (accountID: number): string =>
    `${BASE_URL}player/${accountID},/`;

  export const ship = (shipID: number): string => `${BASE_URL}ship/${shipID},/`;
}
