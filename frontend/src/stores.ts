import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { OptionalBattle, StatsExtra } from "src/lib/types";
import { type Writable, derived, writable } from "svelte/store";
import type { data } from "wailsjs/go/models";

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedConfig = writable({}) as Writable<data.UserConfigV2>;

export const storedAlertPlayers = writable([]) as Writable<data.AlertPlayer[]>;
export const storedLogs = writable([]) as Writable<string[]>;
export const storedInstallPathError = writable("") as Writable<string>;
export const storedTeamThreatLevels = derived(
  [storedBattle, storedConfig],
  ([storedBattle, storedConfig]) =>
    TeamThreatLevel.fromBattle(
      storedBattle,
      storedConfig.stats_pattern as StatsExtra,
    ),
);
