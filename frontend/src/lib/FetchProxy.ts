import {
  storedBattle,
  storedExcludedPlayers,
  storedLogs,
} from "src/stores";
import { Battle, ExcludePlayerIDs } from "wailsjs/go/main/App";
import { data } from "wailsjs/go/models";
import { EventsOn } from "wailsjs/runtime/runtime";

export namespace FetchProxy {
  export const getBattle = async (): Promise<data.Battle> => {
    // Note: 過去のデータが影響してか値が0になってしまうためクリーンする
    storedBattle.set(undefined);
    const ret = await Battle();
    storedBattle.set(ret);
    return ret;
  };

  export const getExcludedPlayers = async (): Promise<number[]> => {
    const ret = await ExcludePlayerIDs();
    storedExcludedPlayers.set(ret);
    return ret;
  };
}

EventsOn("LOG", (log: string) =>
  storedLogs.update((logs) => {
    logs.push(log);
    return logs;
  }),
);
