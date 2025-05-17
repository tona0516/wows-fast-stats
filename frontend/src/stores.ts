import { Summary } from "src/lib/Summary";
import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import type { OptionalBattle, StatsExtra } from "src/lib/types";
import { type Writable, derived, writable } from "svelte/store";
import type { data } from "wailsjs/go/models";

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedExcludedPlayers = writable(new Set()) as Writable<
  Set<number>
>;
export const storedConfig = writable({}) as Writable<data.UserConfigV2>;

export const storedSummary = derived(
  [storedBattle, storedExcludedPlayers, storedConfig],
  ([$storedBattle, $storedExcludedPlayers, $storedConfig]) =>
    Summary.calculate($storedBattle, $storedExcludedPlayers, $storedConfig),
);
export const storedAlertPlayers = writable([]) as Writable<data.AlertPlayer[]>;
export const storedLogs = writable([]) as Writable<string[]>;
export const storedInstallPathError = writable("") as Writable<string>;
export const storedTeamThreatLevels = derived(
  [storedBattle, storedExcludedPlayers, storedConfig],
  ([storedBattle, storedExcludedPlayers, storedConfig]) =>
    TeamThreatLevel.fromBattle(
      storedBattle,
      storedExcludedPlayers,
      storedConfig.stats_pattern as StatsExtra,
    ),
);
