import { ROMAN_NUMERALS } from "./constants";

export const toTierString = (value: number): string => {
  if (value === 11) return "★";
  return ROMAN_NUMERALS[value] ?? "";
};

const formatNumber = (value: number, digit: number): string => {
  return value.toLocaleString("ja-JP", {
    minimumFractionDigits: digit,
    maximumFractionDigits: digit,
  });
};

export const formatWithSuffix = (num: number): string => {
  const G = 1_000_000_000;
  const M = 1_000_000;
  const K = 1_000;

  if (num >= G) {
    return `${formatNumber(num / G, 1)}G`;
  }

  if (num >= M) {
    return `${formatNumber(num / M, 1)}M`;
  }

  if (num >= K) {
    return `${formatNumber(num / K, 1)}K`;
  }

  return formatNumber(num, 0);
};
