import { AppConstants } from "./AppConstants";
import { LocalStorage } from "./LocalStorage";

export namespace AppFunctions {
  export const isLighter = (): boolean => {
    const current = LocalStorage.instance.getTheme();
    return AppConstants.LIGHTER_THEMES.includes(current);
  };
}
