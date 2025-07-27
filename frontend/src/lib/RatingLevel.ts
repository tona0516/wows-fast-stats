import type { data } from "wailsjs/go/models";
import { Color } from "./Color";

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
  color: string;
  pr: number;
  damage: number;
  winRate: number;
}

export const RATING_DEFS: Rating[] = [
  {
    level: "bad",
    color: Color.Rating.getDefault("bad"),
    pr: 0,
    damage: 0,
    winRate: 0,
  },
  {
    level: "below_avg",
    color: Color.Rating.getDefault("below_avg"),
    pr: 750,
    damage: 0.6,
    winRate: 47,
  },
  {
    level: "avg",
    color: Color.Rating.getDefault("avg"),
    pr: 1100,
    damage: 0.8,
    winRate: 50,
  },
  {
    level: "good",
    color: Color.Rating.getDefault("good"),
    pr: 1350,
    damage: 1.0,
    winRate: 52,
  },
  {
    level: "very_good",
    color: Color.Rating.getDefault("very_good"),
    pr: 1550,
    damage: 1.2,
    winRate: 54,
  },
  {
    level: "great",
    color: Color.Rating.getDefault("great"),
    pr: 1750,
    damage: 1.4,
    winRate: 56,
  },
  {
    level: "unicum",
    color: Color.Rating.getDefault("unicum"),
    pr: 2100,
    damage: 1.5,
    winRate: 60,
  },
  {
    level: "super_unicum",
    color: Color.Rating.getDefault("super_unicum"),
    pr: 2450,
    damage: 1.6,
    winRate: 65,
  },
];
