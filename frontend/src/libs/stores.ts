import type { data } from "@wails/go/models";
import { derived, type Writable, writable } from "svelte/store";
import { STATS_KEYS } from "./constants";
import { LocalStorage } from "./LocalStorage";
import type {
  ColumnSetting,
  EditModalParam,
  Optional,
  StatsKey,
  TonakoParam,
} from "./types";
import { getTeamThreatLevels } from "./utils";

export const storedZoomRate = writable(LocalStorage.instance.getZoomRate());
storedZoomRate.subscribe((zoomRate) => {
  LocalStorage.instance.setZoomRate(zoomRate);
});

export const storedStatsExtra = writable(LocalStorage.instance.getStatsExtra());
storedStatsExtra.subscribe((statsExtra) => {
  LocalStorage.instance.setStatsExtra(statsExtra);
});

export const storedPlayerNameColor = writable(
  LocalStorage.instance.getPlayerNameColor(),
);
storedPlayerNameColor.subscribe((color) => {
  LocalStorage.instance.setPlayerNameColor(color);
});

export const storedShowClanNation = writable(
  LocalStorage.instance.getShowClanNation(),
);
storedShowClanNation.subscribe((showClanNation) => {
  LocalStorage.instance.setShowClanNation(showClanNation);
});

export const storedColumnmSettings = writable(
  STATS_KEYS.reduce(
    (acc, key) => {
      acc[key] = LocalStorage.instance.getColumnSetting(key);
      return acc;
    },
    {} as { [key in StatsKey]: ColumnSetting },
  ),
);
storedColumnmSettings.subscribe((settings) => {
  Object.entries(settings).forEach(([key, value]) => {
    LocalStorage.instance.setColumnSetting(key as StatsKey, value);
  });
});

export const storedBattle = writable(undefined) as Writable<
  Optional<data.Battle>
>;
export const storedAlertPlayers = writable([]) as Writable<data.AlertPlayer[]>;
export const storedInstallPathError = writable("") as Writable<string>;
export const storedTeamThreatLevels = derived(
  [storedBattle, storedStatsExtra],
  ([battle, statsExtra]) => {
    return getTeamThreatLevels(battle, statsExtra);
  },
);
export const storedEditAlertPlayer = writable(undefined) as Writable<
  Optional<EditModalParam>
>;
export const storedDeleteAlertPlayer = writable(undefined) as Writable<
  Optional<data.AlertPlayer>
>;
export const storedPlayerDetail = writable(undefined) as Writable<
  Optional<data.Player>
>;
export const storedPlayerShipDetail = writable(undefined) as Writable<
  Optional<data.Player>
>;

export const storedTonako = writable(undefined) as Writable<
  Optional<TonakoParam>
>;

export const storedToastText = writable("");
export function showToast(text: string, intervalSeconds = 3) {
  storedToastText.set(text);
  setTimeout(() => {
    storedToastText.set("");
  }, intervalSeconds * 1000);
}
