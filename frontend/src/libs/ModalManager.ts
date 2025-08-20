import type { data } from "@wails/go/models";
import { DEFAULT_ALERT_PLAYER } from "./constants";
import {
  storedEditAlertPlayer,
  storedPlayerDetail,
  storedDeleteAlertPlayer as storedRemoveAlertPlayer,
  storedPlayerShipDetail as storedShipDetail,
} from "./stores";

export class ModalManager {
  private static _instance: ModalManager;

  private constructor() {}

  public static get instance(): ModalManager {
    if (!ModalManager._instance) {
      ModalManager._instance = new ModalManager();
    }
    return ModalManager._instance;
  }

  openForCreate() {
    storedEditAlertPlayer.set({
      mode: "create",
      form: structuredClone(DEFAULT_ALERT_PLAYER),
    });
  }

  openForSpecify(accountID: number, name: string) {
    const defaultValue = structuredClone(DEFAULT_ALERT_PLAYER);
    storedEditAlertPlayer.set({
      mode: "specify",
      form: {
        account_id: accountID,
        name: name,
        pattern: defaultValue.pattern,
        message: defaultValue.message,
      } as data.AlertPlayer,
    });
  }

  openForEdit(ap: data.AlertPlayer) {
    storedEditAlertPlayer.set({
      mode: "edit",
      form: ap,
    });
  }

  closeForEdit() {
    storedEditAlertPlayer.set(undefined);
  }

  openForRemove(ap: data.AlertPlayer) {
    storedRemoveAlertPlayer.set(ap);
  }

  closeForDelete() {
    storedRemoveAlertPlayer.set(undefined);
  }

  openForPlayerDetail(player: data.Player) {
    storedPlayerDetail.set(player);
  }

  closeForPlayerDetail() {
    storedPlayerDetail.set(undefined);
  }

  openForShipDetail(player: data.Player) {
    storedShipDetail.set(player);
  }

  closeShipDetail() {
    storedShipDetail.set(undefined);
  }
}
