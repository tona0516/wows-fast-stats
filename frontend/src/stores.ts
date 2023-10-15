import { Summary } from "src/lib/Summary";
import { type OptionalBattle } from "src/lib/types";
import { derived, writable, type Writable } from "svelte/store";
import { domain, vo } from "wailsjs/go/models";

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedExcludedPlayers = writable([]) as Writable<number[]>;
export const storedConfig = writable({}) as Writable<domain.UserConfig>;

export const storedSummary = derived(
  [storedBattle, storedExcludedPlayers, storedConfig],
  ([$storedBattle, $storedExcludePlayerIDs, $storedUserConfig]) =>
    Summary.calculate(
      $storedBattle,
      $storedExcludePlayerIDs,
      $storedUserConfig,
    ),
);
export const storedAlertPlayers = writable([]) as Writable<
  domain.AlertPlayer[]
>;
export const storedLogs = writable([]) as Writable<string[]>;
export const storedRequiredConfigError = writable(
  {},
) as Writable<vo.RequiredConfigError>;
