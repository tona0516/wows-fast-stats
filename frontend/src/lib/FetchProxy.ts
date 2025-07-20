import { storedBattle } from "src/stores";
import { Battle } from "wailsjs/go/main/App";
import type { data } from "wailsjs/go/models";

export namespace FetchProxy {
  export const getBattle = async (): Promise<data.Battle> => {
    // Note: 過去のデータが影響してか値が0になってしまうためクリーンする
    storedBattle.set(undefined);
    const cache = localStorage.getItem("cache");

    // TODO: あとで消す
    if (cache) {
      const cachedBattle = JSON.parse(cache) as data.Battle;
      storedBattle.set(cachedBattle);
      return cachedBattle;
    }

    const ret = await Battle();
    storedBattle.set(ret);
    localStorage.setItem("cache", JSON.stringify(ret));
    return ret;
  };
}
