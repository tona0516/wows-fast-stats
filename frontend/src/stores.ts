import { Page, type OptionalBattle, type OptionalSummary } from "src/lib/types";
import { writable, type Writable } from "svelte/store";
import { domain } from "wailsjs/go/models";

export const storedBattle = writable(undefined) as Writable<OptionalBattle>;
export const storedSummary = writable(undefined) as Writable<OptionalSummary>;
export const storedExcludePlayerIDs = writable([]) as Writable<number[]>;
export const storedCurrentPage = writable(Page.MAIN) as Writable<Page>;
export const storedUserConfig = writable({}) as Writable<domain.UserConfig>;
export const storedDefaultUserConfig = writable(
  {},
) as Writable<domain.UserConfig>;
export const storedAlertPlayers = writable([]) as Writable<
  domain.AlertPlayer[]
>;
export const storedConfirmMessage = writable("") as Writable<string>;
export const storedLogs = writable([]) as Writable<string[]>;
