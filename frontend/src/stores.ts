import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { ColumnSetting, OptionalBattle, StatsKey } from "src/lib/types";
import { derived, type Writable, writable } from "svelte/store";
import type { data } from "wailsjs/go/models";
import type { Tonako } from "./component/stats/internal/Tonako";
import { AppConstants } from "./lib/AppConstants";
import { LocalStorage } from "./lib/LocalStorage";

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
  AppConstants.STATS_KEYS.reduce(
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

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedAlertPlayers = writable([]) as Writable<data.AlertPlayer[]>;
export const storedInstallPathError = writable("") as Writable<string>;
export const storedTeamThreatLevels = derived(
  [storedBattle, storedStatsExtra],
  ([battle, statsExtra]) => {
    return TeamThreatLevel.fromBattle(battle, statsExtra);
  },
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
