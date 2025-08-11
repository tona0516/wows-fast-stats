import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { OptionalBattle, StatsExtra } from "src/lib/types";
import { derived, type Writable, writable } from "svelte/store";
import type { data } from "wailsjs/go/models";
import type { Tonako } from "./component/stats/internal/Tonako";

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
export const storedTonako = writable(undefined) as Writable<
  { message: string; isLoading: boolean; tonako: Tonako } | undefined
>;

const DEFAULT_ALERT_PLAYER = {
  account_id: 0,
  name: "",
  pattern: "bi-check-circle-fill",
  message: "",
} as data.AlertPlayer;
type EditModalMode = "create" | "specify" | "edit";
export const storedEditAlertPlayer = writable(undefined) as Writable<
  | {
      mode: EditModalMode;
      form: data.AlertPlayer;
    }
  | undefined
>;
export class EditAlertPlayerModal {
  private constructor() {}

  static openForCreate() {
    storedEditAlertPlayer.set({
      mode: "create",
      form: structuredClone(DEFAULT_ALERT_PLAYER),
    });
  }

  static openForSpecify(accountID: number, name: string) {
    const defaultValue = structuredClone(DEFAULT_ALERT_PLAYER);
    storedEditAlertPlayer.set({
      mode: "specify",
      form: {
        account_id: accountID,
        name: name,
        pattern: defaultValue.pattern,
        message: defaultValue.message,
      } as data.AlertPlayer,
    });
  }

  static openForEdit(ap: data.AlertPlayer) {
    storedEditAlertPlayer.set({
      mode: "edit",
      form: ap,
    });
  }

  static close() {
    storedEditAlertPlayer.set(undefined);
  }
}

export const storedDeleteAlertPlayer = writable(undefined) as Writable<
  data.AlertPlayer | undefined
>;
export class DeleteAlertPlayerModal {
  private constructor() {}

  static open(ap: data.AlertPlayer) {
    storedDeleteAlertPlayer.set(ap);
  }

  static close() {
    storedDeleteAlertPlayer.set(undefined);
  }
}

export const storedPlayerDetail = writable(undefined) as Writable<
  data.Player | undefined
>;
export class PlayerDetailModal {
  private constructor() {}

  static open(player: data.Player) {
    storedPlayerDetail.set(player);
  }

  static close() {
    storedPlayerDetail.set(undefined);
  }
}

export const storedPlayerDetailForShip = writable(undefined) as Writable<
  data.Player | undefined
>;
export class ShipDetailModal {
  private constructor() {}

  static open(player: data.Player) {
    storedPlayerDetailForShip.set(player);
  }

  static close() {
    storedPlayerDetailForShip.set(undefined);
  }
}

export const storedToastText = writable("");
export function showToast(text: string, intervalSeconds = 3) {
  storedToastText.set(text);
  setTimeout(() => {
    storedToastText.set("");
  }, intervalSeconds * 1000);
}
