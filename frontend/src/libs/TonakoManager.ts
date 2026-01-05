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

  setNeedInitialSettingState() {
    storedTonako.set({
      message: "設定画面から初期設定をおこなってください",
      isLoading: false,
      tonako: Tonako.Pointing,
    });
  }

  setStartBattleState() {
    storedTonako.set({
      message: "戦闘データを読み込み中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  }

  setPollingStartState() {
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
