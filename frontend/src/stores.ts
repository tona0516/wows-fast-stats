import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { OptionalBattle, PlayerDetail, StatsExtra } from "src/lib/types";
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
export const storedIsShowPlayerDetailModal = writable(
  false,
) as Writable<boolean>;
export const storedAlertPlayerForm = writable({
  account_id: 0,
  name: "",
  pattern: "bi-check-circle-fill",
  message: "",
}) as Writable<data.AlertPlayer>;
export const storedPlayerDetail = writable({
  id: 0,
  name: "",
  clan: undefined,
}) as Writable<PlayerDetail>;
export const storedIsEditAlertPlayer = writable(false) as Writable<boolean>;

export const storedToastText = writable("");

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

export const showPlayerDetailModal = (player: PlayerDetail) => {
  storedIsShowPlayerDetailModal.set(true);
  storedPlayerDetail.set(player);
};

export const closeModal = () => {
  storedIsShowUpdateAlertPlayerModal.set(false);
  storedIsShowDeleteAlertPlayerModal.set(false);
  storedIsShowPlayerDetailModal.set(false);
  storedAlertPlayerForm.set({
    account_id: 0,
    name: "",
    pattern: "",
    message: "",
  });
  storedPlayerDetail.set({
    id: 0,
    name: "",
    clan: undefined,
  });
  storedIsEditAlertPlayer.set(false);
};

export function showToast(text: string, intervalSeconds = 3) {
  storedToastText.set(text);
  setTimeout(() => {
    storedToastText.set("");
  }, intervalSeconds * 1000);
}
