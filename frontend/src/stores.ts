import { writable, type Writable } from "svelte/store";
import type { vo } from "../wailsjs/go/models";
import type { Page } from "./enums";
import type { SummaryResult } from "./util";

export const storedBattle = writable(undefined) as Writable<vo.Battle>;
export const storedSummaryResult = writable(
  undefined
) as Writable<SummaryResult>;
export const storedExcludePlayerIDs = writable([]) as Writable<number[]>;
export const storedCurrentPage = writable("main") as Writable<Page>;
export const storedUserConfig = writable({}) as Writable<vo.UserConfig>;
export const storedAlertPlayers = writable([]) as Writable<vo.AlertPlayer[]>;
export const storedIsFirstScreenshot = writable(true) as Writable<boolean>;
