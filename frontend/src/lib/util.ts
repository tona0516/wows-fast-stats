import type { ColumnSetting } from "src/lib/ColumnSetting";
import type { DigitKey, OverallKey, ShipKey, ShipType } from "src/lib/types";
import { data } from "wailsjs/go/models";

const ROMAN_NUMERALS: { [key: number]: string } = {
  1: "I",
  2: "II",
  3: "III",
  4: "IV",
  5: "V",
  6: "VI",
  7: "VII",
  8: "VIII",
  9: "IX",
  10: "X",
};

export const toPlayerStats = (
  player: data.Player,
  statsPattern: string,
): data.PlayerStats => {
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
  return ROMAN_NUMERALS[value];
};

export const isShipType = (type: string): type is ShipType => {
  return Object.keys(new data.ShipTypeGroup()).includes(type);
};

export const isDigitKey = (key: string): key is DigitKey => {
  return Object.keys(new data.UCDigit()).includes(key);
};

export const isShipKey = (key: string): key is ShipKey => {
  return Object.keys(new data.UCDisplayShip()).includes(key);
};

export const isOverallKey = (key: string): key is OverallKey => {
  return Object.keys(new data.UCDisplayOverall()).includes(key);
};

export const deriveColumnSetting = (
  config: data.UserConfigV2,
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
  config: data.UserConfigV2,
): ColumnSetting[] => {
  const shipKeys = Object.keys(config.display.ship);
  const overallKeys = Object.keys(config.display.overall);
  const allKeys = [...new Set([...shipKeys, ...overallKeys])];

  return allKeys.map((key) => deriveColumnSetting(config, key));
};
