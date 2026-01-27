import type { core } from "@wails/go/models";
import { type Writable, writable } from "svelte/store";
import type { DisplayPref } from "./DisplayPref";
import type { Optional, TonakoParam } from "./types";

export const storedGameClientPath = writable("") as Writable<string>;
export const storedDisplayPref = writable() as Writable<DisplayPref>;
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
