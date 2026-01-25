import type { core } from "@wails/go/models";
import { type Writable, writable } from "svelte/store";
import type { Optional, TonakoParam } from "./types";

export const storedPref = writable() as Writable<core.Pref>;
export const storedBattle = writable(undefined) as Writable<
  Optional<core.Battle>
>;
export const storedPlayerDetail = writable(undefined) as Writable<
  Optional<core.Player>
>;
export const storedPlayerShipDetail = writable(undefined) as Writable<
  Optional<core.Player>
>;
export const storedTonako = writable(undefined) as Writable<
  Optional<TonakoParam>
>;

export const storedToastText = writable("");
export function showToast(text: string, intervalSeconds = 5) {
  storedToastText.set(text);
  setTimeout(() => {
    storedToastText.set("");
  }, intervalSeconds * 1000);
}

export const storedGameClientPathError = writable("") as Writable<string>;
