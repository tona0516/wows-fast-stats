import {
  UpdateBasicColumnSetting,
  UpdateOptionalSetting,
  UpdateRequiredSetting,
  UpdateStatsColumnSettings,
} from "@wails/go/main/App";
import type { data } from "@wails/go/models";
import { derived, type Writable, writable } from "svelte/store";
import type {
  EditModalParam,
  Optional,
  StatsExtra,
  TonakoParam,
} from "./types";
import { getTeamThreatLevels } from "./utils";

export const storedRequiredSetting =
  writable() as Writable<data.RequiredSetting>;
storedRequiredSetting.subscribe(async (value) => {
  if (!value) return;
  await UpdateRequiredSetting(value);
});

export const storedOptionalSetting =
  writable() as Writable<data.OptionalSetting>;
storedOptionalSetting.subscribe(async (value) => {
  if (!value) return;
  document.body.style.zoom = `${value.zoom_rate}%`;
  await UpdateOptionalSetting(value);
});

export const storedBasicColumnSetting =
  writable() as Writable<data.BasicColumnSetting>;
storedBasicColumnSetting.subscribe(async (value) => {
  if (!value) return;
  await UpdateBasicColumnSetting(value);
});

export const storedStatsColumnSettings =
  writable() as Writable<data.StatsColumnSettings>;
storedStatsColumnSettings.subscribe(async (value) => {
  if (!value) return;
  await UpdateStatsColumnSettings(value);
});

export const storedBattle = writable(undefined) as Writable<
  Optional<data.Battle>
>;
export const storedAlertPlayers = writable([]) as Writable<data.AlertPlayer[]>;
export const storedInstallPathError = writable("") as Writable<string>;
export const storedTeamThreatLevels = derived(
  [storedBattle, storedOptionalSetting],
  ([battle, optionalSetting]) => {
    return getTeamThreatLevels(
      battle,
      optionalSetting.stats_extra as StatsExtra,
    );
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
