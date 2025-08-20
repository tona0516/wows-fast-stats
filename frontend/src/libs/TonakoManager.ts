import { storedTonako } from "./stores";
import { Tonako } from "./Tonako";

export class TonakoManager {
  private static _instance: TonakoManager;

  private constructor() {}

  public static get instance(): TonakoManager {
    if (!TonakoManager._instance) {
      TonakoManager._instance = new TonakoManager();
    }
    return TonakoManager._instance;
  }

  setStartBattleState() {
    storedTonako.set({
      message: "戦闘データを読み込み中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  }

  setEndBattleState() {
    storedTonako.set({
      message: "戦闘開始時に自動的にリロードします",
      isLoading: false,
      tonako: Tonako.Standby,
    });
  }

  setBattleErrorState(message: string) {
    storedTonako.set({
      message: message,
      isLoading: false,
      tonako: Tonako.Sorry,
    });
  }

  setFetchOtherDataState() {
    storedTonako.set({
      message: "艦・マップ情報を取得中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  }

  setFetchPlayerDataState() {
    storedTonako.set({
      message: "プレイヤー情報を取得中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  }

  setHidden() {
    storedTonako.set(undefined);
  }
}
