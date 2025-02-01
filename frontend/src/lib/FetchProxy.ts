import { storedBattle } from "src/stores";
import { Battle } from "wailsjs/go/main/App";
import { model } from "wailsjs/go/models";

export namespace FetchProxy {
  export const getBattle = async (): Promise<model.Battle> => {
    // Note: 過去のデータが影響してか値が0になってしまうためクリーンする
    storedBattle.set(undefined);
    const ret = await Battle();
    storedBattle.set(ret);
    return ret;
  };
}
