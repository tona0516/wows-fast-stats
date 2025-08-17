import { AppConst } from "./AppConst";
import { LocalStorage } from "./LocalStorage";

export namespace AppFunc {
  export const isLighter = (): boolean => {
    const current = LocalStorage.instance.getTheme();
    return AppConst.LIGHTER_THEMES.includes(current);
  };

  export const clanNumbersURL = (clanID: number): string =>
    `${AppConst.NUMBERS_URL}clan/${clanID},/`;

  export const playerNumbersURL = (
    accountID: number,
    accountName: string,
  ): string => `${AppConst.NUMBERS_URL}player/${accountID},${accountName}/`;

  export const shipNumbersURL = (shipID: number): string =>
    `${AppConst.NUMBERS_URL}ship/${shipID},/`;
}
