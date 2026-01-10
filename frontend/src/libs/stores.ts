import { SaveUserConfig } from "@wails/go/main/App";
import type { data } from "@wails/go/models";
import { type Writable, writable } from "svelte/store";
import type { Optional, TonakoParam } from "./types";

export const storedUserConfig = writable() as Writable<data.UserConfig>;
storedUserConfig.subscribe(async (value) => {
  if (!value) return;
  await SaveUserConfig(value);
});

export const storedBattle = writable(undefined) as Writable<
  Optional<data.Battle>
>;
export const storedPlayerDetail = writable(undefined) as Writable<
  Optional<data.Player>
>;
export const storedPlayerShipDetail = writable(undefined) as Writable<
  Optional<data.Player>
>;
export const storedTonako = writable(undefined) as Writable<
  Optional<TonakoParam>
>;

export const storedToastText = writable("");
export function showToast(text: string, intervalSeconds = 3) {
  storedToastText.set(text);
  setTimeout(() => {
    storedToastText.set("");
  }, intervalSeconds * 1000);
}
