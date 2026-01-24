import { storedTonako } from "./stores";
import { Tonako } from "./Tonako";

export class TonakoManager {
  private static instance: TonakoManager;

  private constructor() {}

  public static get getInstance(): TonakoManager {
    if (!TonakoManager.instance) {
      TonakoManager.instance = new TonakoManager();
    }
    return TonakoManager.instance;
  }

  setPromoteState(message: string) {
    storedTonako.set({
      message: message,
      isLoading: false,
      tonako: Tonako.Pointing,
    });
  }

  setStandbyState(message: string) {
    storedTonako.set({
      message: message,
      isLoading: false,
      tonako: Tonako.Standby,
    });
  }

  setLoadingState(message: string) {
    storedTonako.set({
      message: message,
      isLoading: true,
      tonako: Tonako.Standby,
    });
  }

  setErrorState(message: string) {
    storedTonako.set({
      message: message,
      isLoading: false,
      tonako: Tonako.Sorry,
    });
  }

  setHidden() {
    storedTonako.set(undefined);
  }
}
