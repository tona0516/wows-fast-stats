import type { data } from "wailsjs/go/models";

export type RatingLevel = Readonly<keyof data.UCSkillColorCode>;

export class RatingfGenerator {
  private constructor() {}

  static fromPR(value: number): Rating | undefined {
    return RATING_DEFS.findLast((it) => value >= it.pr);
  }

  static fromWinRate(value: number): Rating | undefined {
    return RATING_DEFS.findLast((it) => value >= it.winRate);
  }

  static fromDamage(value: number, expected: number): Rating | undefined {
    if (expected === 0) {
      return undefined;
    }

    const ratio = value / expected;
    return RATING_DEFS.findLast((it) => ratio >= it.damage);
  }
}

export interface Rating {
  level: RatingLevel;
  textColor: string;
  bgColor: string;
  pr: number;
  damage: number;
  winRate: number;
}

export const RATING_DEFS: Rating[] = [
  {
    level: "bad",
    textColor: "#ffffff",
    bgColor: "#FE0E00",
    pr: 0,
    damage: 0,
    winRate: 0,
  },
  {
    level: "below_avg",
    textColor: "#000000",
    bgColor: "#FE7903",
    pr: 750,
    damage: 0.6,
    winRate: 47,
  },
  {
    level: "avg",
    textColor: "#000000",
    bgColor: "#FFC71F",
    pr: 1100,
    damage: 0.8,
    winRate: 50,
  },
  {
    level: "good",
    textColor: "#ffffff",
    bgColor: "#44B300",
    pr: 1350,
    damage: 1.0,
    winRate: 52,
  },
  {
    level: "very_good",
    textColor: "#ffffff",
    bgColor: "#318000",
    pr: 1550,
    damage: 1.2,
    winRate: 54,
  },
  {
    level: "great",
    textColor: "#ffffff",
    bgColor: "#02C9B3",
    pr: 1750,
    damage: 1.4,
    winRate: 56,
  },
  {
    level: "unicum",
    textColor: "#ffffff",
    bgColor: "#D042F3",
    pr: 2100,
    damage: 1.5,
    winRate: 60,
  },
  {
    level: "super_unicum",
    textColor: "#ffffff",
    bgColor: "#A00DC5",
    pr: 2450,
    damage: 1.6,
    winRate: 65,
  },
];
