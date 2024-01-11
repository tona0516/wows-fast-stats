import {
  storedAlertPlayers,
  storedBattle,
  storedConfig,
  storedExcludedPlayers,
  storedLogs,
  storedRequiredConfigError,
} from "src/stores";
import {
  AlertPlayers,
  ApplyRequiredUserConfig,
  Battle,
  ExcludePlayerIDs,
  UserConfig,
  ValidateRequiredConfig,
} from "wailsjs/go/main/App";
import { model } from "wailsjs/go/models";
import { EventsOn } from "wailsjs/runtime/runtime";

export namespace FetchProxy {
  export const getBattle = async (): Promise<model.Battle> => {
    const ret = await Battle();
    storedBattle.set(ret);
    return ret;
  };

  export const getConfig = async (): Promise<model.UserConfig> => {
    const ret = await UserConfig();
    storedConfig.set(ret);
    return ret;
  };

  export const validateRequiredConfig = async (
    installPath: string,
    appid: string,
  ): Promise<model.RequiredConfigError> => {
    const ret = await ValidateRequiredConfig(installPath, appid);
    storedRequiredConfigError.set(ret);
    return ret;
  };

  export const applyRequiredConfig = async (
    installPath: string,
    appid: string,
  ): Promise<model.RequiredConfigError> => {
    const ret = await ApplyRequiredUserConfig(installPath, appid);
    storedRequiredConfigError.set(ret);
    return ret;
  };

  export const getAlertPlayers = async (): Promise<model.AlertPlayer[]> => {
    const ret = await AlertPlayers();
    storedAlertPlayers.set(ret);
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
