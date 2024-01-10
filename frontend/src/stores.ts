import { Summary } from "src/lib/Summary";
import { type OptionalBattle } from "src/lib/types";
import { derived, writable, type Writable } from "svelte/store";
import { domain } from "wailsjs/go/models";

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedExcludedPlayers = writable([]) as Writable<number[]>;
export const storedConfig = writable({}) as Writable<domain.UserConfig>;

export const storedSummary = derived(
  [storedBattle, storedExcludedPlayers, storedConfig],
  ([$storedBattle, $storedExcludedPlayers, $storedConfig]) =>
    Summary.calculate($storedBattle, $storedExcludedPlayers, $storedConfig),
);
export const storedAlertPlayers = writable([]) as Writable<
  domain.AlertPlayer[]
>;
export const storedLogs = writable([]) as Writable<string[]>;
export const storedRequiredConfigError = writable(
  {},
) as Writable<domain.RequiredConfigError>;
