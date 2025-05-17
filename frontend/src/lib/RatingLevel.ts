import type { data } from "wailsjs/go/models";

export type RatingLevel = Readonly<keyof data.UCSkillColorCode>;

export class RatingInfo {
  constructor(
    readonly level: RatingLevel,
    readonly textColorCode: string,
  ) {}

  static fromPR(
    value: number,
    colorCode: data.UCSkillColorCode,
  ): RatingInfo | undefined {
    const rating = RATING_DEFS.findLast((it) => value >= it.pr);
    return !rating
      ? undefined
      : new RatingInfo(rating.level, colorCode[rating.level]);
  }

  static fromWinRate(
    value: number,
    colorCode: data.UCSkillColorCode,
  ): RatingInfo | undefined {
    const rating = RATING_DEFS.findLast((it) => value >= it.winRate);
    return !rating
      ? undefined
      : new RatingInfo(rating.level, colorCode[rating.level]);
  }

  static fromDamage(
    value: number,
    expected: number,
    colorCode: data.UCSkillColorCode,
  ): RatingInfo | undefined {
    if (expected === 0) {
      return undefined;
    }

    const ratio = value / expected;
    const rating = RATING_DEFS.findLast((it) => ratio >= it.damage);
    return !rating
      ? undefined
      : new RatingInfo(rating?.level, colorCode[rating.level]);
  }
}

export type RatingDef = {
  level: RatingLevel;
  pr: number;
  damage: number;
  winRate: number;
};
export const RATING_DEFS: RatingDef[] = [
  {
    level: "bad",
    pr: 0,
    damage: 0,
    winRate: 0,
  },
  {
    level: "below_avg",
    pr: 750,
    damage: 0.6,
    winRate: 47,
  },
  {
    level: "avg",
    pr: 1100,
    damage: 0.8,
    winRate: 50,
  },
  {
    level: "good",
    pr: 1350,
    damage: 1.0,
    winRate: 52,
  },
  {
    level: "very_good",
    pr: 1550,
    damage: 1.2,
    winRate: 54,
  },
  {
    level: "great",
    pr: 1750,
    damage: 1.4,
    winRate: 56,
  },
  {
    level: "unicum",
    pr: 2100,
    damage: 1.5,
    winRate: 60,
  },
  {
    level: "super_unicum",
    pr: 2450,
    damage: 1.6,
    winRate: 65,
  },
];
