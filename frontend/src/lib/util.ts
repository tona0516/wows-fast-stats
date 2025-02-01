import type { ColumnSetting } from "src/lib/ColumnSetting";
import type { DigitKey, OverallKey, ShipKey, ShipType } from "src/lib/types";
import { model } from "wailsjs/go/models";

const ROMAN_NUMERALS: { decimal: number; numeral: string }[] = [
  { decimal: 10, numeral: "X" },
  { decimal: 9, numeral: "IX" },
  { decimal: 5, numeral: "V" },
  { decimal: 4, numeral: "IV" },
  { decimal: 1, numeral: "I" },
];

export const toPlayerStats = (
  player: model.Player,
  statsPattern: string,
): model.PlayerStats => {
  switch (statsPattern) {
    case "pvp_solo":
      return player.pvp_solo;
    case "pvp_all":
      return player.pvp_all;
    case "rank_solo":
      return player.rank_solo;
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
  return Object.keys(new model.ShipTypeGroup()).includes(type);
};

export const isDigitKey = (key: string): key is DigitKey => {
  return Object.keys(new model.UCDigit()).includes(key);
};

export const isShipKey = (key: string): key is ShipKey => {
  return Object.keys(new model.UCDisplayShip()).includes(key);
};

export const isOverallKey = (key: string): key is OverallKey => {
  return Object.keys(new model.UCDisplayOverall()).includes(key);
};

export const deriveColumnSetting = (
  config: model.UserConfigV2,
  key: string,
): ColumnSetting => {
  const shipKey = isShipKey(key) ? key : undefined;
  const overallKey = isOverallKey(key) ? key : undefined;
  const digitKey = isDigitKey(key) ? key : undefined;

  return {
    key: key,
    ship: {
      key: shipKey,
      value: shipKey ? config.display.ship[shipKey] : false,
    },
    overall: {
      key: overallKey,
      value: overallKey ? config.display.overall[overallKey] : false,
    },
    digit: {
      key: digitKey,
      value: digitKey ? config.digit[digitKey] : 0,
    },
  };
};

export const deriveColumnSettings = (
  config: model.UserConfigV2,
): ColumnSetting[] => {
  const shipKeys = Object.keys(config.display.ship);
  const overallKeys = Object.keys(config.display.overall);
  const allKeys = [...new Set([...shipKeys, ...overallKeys])];

  return allKeys.map((key) => deriveColumnSetting(config, key));
};
