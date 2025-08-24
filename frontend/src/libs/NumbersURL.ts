export namespace NumbersURL {
  const NUMBERS_URL = "https://asia.wows-numbers.com/";

  export const getClan = (clanID: number): string =>
    `${NUMBERS_URL}clan/${clanID},/`;

  export const getPlayer = (accountID: number, accountName: string): string =>
    `${NUMBERS_URL}player/${accountID},${accountName}/`;

  export const getShip = (shipID: number): string =>
    `${NUMBERS_URL}ship/${shipID},/`;
}
