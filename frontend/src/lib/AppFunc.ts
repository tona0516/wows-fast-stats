import { AppConst } from "./AppConst";
import { LocalStorage } from "./LocalStorage";

export namespace AppFunc {
  export const isLighter = (): boolean => {
    const current = LocalStorage.instance.getTheme();
    return AppConst.LIGHTER_THEMES.includes(current);
  };
}
