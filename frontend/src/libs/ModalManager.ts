import type { data } from "@wails/go/models";
import {
  storedPlayerDetail,
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
