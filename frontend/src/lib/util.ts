import type { DigitKey, OverallKey, ShipKey, ShipType } from "src/lib/types";
import { domain } from "wailsjs/go/models";

const ROMAN_NUMERALS: { decimal: number; numeral: string }[] = [
  { decimal: 10, numeral: "X" },
  { decimal: 9, numeral: "IX" },
  { decimal: 5, numeral: "V" },
  { decimal: 4, numeral: "IV" },
  { decimal: 1, numeral: "I" },
];

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

  let romanized = "";
  for (const { decimal, numeral } of ROMAN_NUMERALS) {
    while (decimal <= value) {
      romanized += numeral;
      value -= decimal;
    }
  }

  return romanized;
};

export const isShipType = (type: string): type is ShipType => {
  return Object.keys(new domain.ShipTypeGroup()).includes(type);
};

export const isDigitKey = (key: string): key is DigitKey => {
  return Object.keys(new domain.CustomDigit()).includes(key);
};

export const isShipKey = (key: string): key is ShipKey => {
  return Object.keys(new domain.Ship()).includes(key);
};

export const isOverallKey = (key: string): key is OverallKey => {
  return Object.keys(new domain.Overall()).includes(key);
};
