import type { data, domain } from "@wails/go/models";
import { DEFAULT_BLACK_LIST_ITEM } from "./constants";
import {
  storedEditBlackListItem,
  storedPlayerDetail,
  storedRemoveBlackListItem,
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
    storedEditBlackListItem.set({
      mode: "create",
      form: structuredClone(DEFAULT_BLACK_LIST_ITEM),
    });
  }

  openForSpecify(accountID: number, name: string) {
    const defaultValue = structuredClone(DEFAULT_BLACK_LIST_ITEM);
    storedEditBlackListItem.set({
      mode: "specify",
      form: {
        account_id: accountID,
        name: name,
        pattern: defaultValue.pattern,
        message: defaultValue.message,
        created_at: defaultValue.created_at,
      } as domain.BlackListItem,
    });
  }

  openForEdit(item: domain.BlackListItem) {
    storedEditBlackListItem.set({
      mode: "edit",
      form: item,
    });
  }

  closeForEdit() {
    storedEditBlackListItem.set(undefined);
  }

  openForRemove(item: domain.BlackListItem) {
    storedRemoveBlackListItem.set(item);
  }

  closeForDelete() {
    storedRemoveBlackListItem.set(undefined);
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
