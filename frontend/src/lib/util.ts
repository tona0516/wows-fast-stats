// @ts-ignore
import { domain } from "wailsjs/go/models";

export const toPlayerStats = (
  player: domain.Player,
  statsPattern: string,
): domain.PlayerStats => {
  switch (statsPattern) {
    case "pvp_solo":
      return player.pvp_solo;
    case "pvp_all":
      return player.pvp_all;
    default:
      throw Error(`unexpeted error: statsPattern: ${statsPattern}`);
  }
};

export const tierString = (value: number): string => {
  if (value === 11) return "★";

  const decimal = [10, 9, 5, 4, 1];
  const romanNumeral = ["X", "IX", "V", "IV", "I"];

  let romanized = "";

  for (var i = 0; i < decimal.length; i++) {
    while (decimal[i] <= value) {
      romanized += romanNumeral[i];
      value -= decimal[i];
    }
  }

  return romanized;
};
