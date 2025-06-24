import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { OptionalBattle, StatsExtra } from "src/lib/types";
import { type Writable, derived, writable } from "svelte/store";
import type { data } from "wailsjs/go/models";

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedConfig = writable({}) as Writable<data.UserConfigV2>;

export const storedAlertPlayers = writable([]) as Writable<data.AlertPlayer[]>;
export const storedInstallPathError = writable("") as Writable<string>;
export const storedTeamThreatLevels = derived(
  [storedBattle, storedConfig],
  ([storedBattle, storedConfig]) =>
    TeamThreatLevel.fromBattle(
      storedBattle,
      storedConfig.stats_pattern as StatsExtra,
    ),
);

export const storedIsShowUpdateAlertPlayerModal = writable(
  false,
) as Writable<boolean>;
export const storedIsShowDeleteAlertPlayerModal = writable(
  false,
) as Writable<boolean>;
export const storedAlertPlayerForm = writable({
  account_id: 0,
  name: "",
  pattern: "bi-check-circle-fill",
  message: "",
}) as Writable<data.AlertPlayer>;
export const storedIsEditAlertPlayer = writable(false) as Writable<boolean>;

export const showAddAlertPlayerModal = () => {
  storedIsShowUpdateAlertPlayerModal.set(true);
  storedIsEditAlertPlayer.set(false);
};

export const showUpdateAlertPlayerModal = (
  accountID: number,
  name: string,
  pattern?: string,
  message?: string,
) => {
  storedIsShowUpdateAlertPlayerModal.set(true);
  storedIsEditAlertPlayer.set(true);
  storedAlertPlayerForm.set({
    account_id: accountID,
    name: name,
    pattern: pattern || "bi-check-circle-fill",
    message: message || "",
  });
};

export const showDeleteAlertPlayerModal = (accountID: number) => {
  storedIsShowDeleteAlertPlayerModal.set(true);
  storedAlertPlayerForm.set({
    account_id: accountID,
    name: "",
    pattern: "",
    message: "",
  });
};

export const closeAlertPlayerModal = () => {
  storedIsShowUpdateAlertPlayerModal.set(false);
  storedIsShowDeleteAlertPlayerModal.set(false);
  storedAlertPlayerForm.set({
    account_id: 0,
    name: "",
    pattern: "",
    message: "",
  });
  storedIsEditAlertPlayer.set(false);
};
