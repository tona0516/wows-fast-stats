import type { ShipType } from "src/lib/types";
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

export const tierString = (value: number): string => {
  if (value === 11) return "★";
  return ROMAN_NUMERALS[value] ?? "";
};

export const isShipType = (type: string): type is ShipType => {
  return Object.keys(new data.ShipTypeGroup()).includes(type);
};

const G = 1_000_000_000;
const M = 1_000_000;
const K = 1_000;

export const formatWithSuffix = (num: number): string => {
  if (num >= G) {
    return `${(num / G).format(1)}G`;
  }

  if (num >= M) {
    return `${(num / M).format(1)}M`;
  }

  if (num >= K) {
    return `${(num / K).format(1)}K`;
  }

  return num.format(1);
};

declare global {
  interface Number {
    format(digit: number): string;
  }
}

Number.prototype.format = function (digit: number): string {
  return this.toLocaleString("ja-JP", {
    minimumFractionDigits: digit,
    maximumFractionDigits: digit,
  });
};
