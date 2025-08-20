import { ColorCode } from "./ColorCode";
import type { Optional } from "./types";

type RawRating =
  | "bad"
  | "below_avg"
  | "avg"
  | "good"
  | "very_good"
  | "great"
  | "unicum"
  | "super_unicum";

type RatingThreshold = {
  readonly raw: RawRating;
  readonly pr: number;
  readonly shipDamageRatio: number;
  readonly winRate: number;
};

const THRESHOLDS: RatingThreshold[] = [
  {
    raw: "bad",
    pr: 0,
    shipDamageRatio: 0,
    winRate: 0,
  },
  {
    raw: "below_avg",
    pr: 750,
    shipDamageRatio: 0.6,
    winRate: 47,
  },
  {
    raw: "avg",
    pr: 1100,
    shipDamageRatio: 0.8,
    winRate: 50,
  },
  {
    raw: "good",
    pr: 1350,
    shipDamageRatio: 1.0,
    winRate: 52,
  },
  {
    raw: "very_good",
    pr: 1550,
    shipDamageRatio: 1.2,
    winRate: 54,
  },
  {
    raw: "great",
    pr: 1750,
    shipDamageRatio: 1.4,
    winRate: 56,
  },
  {
    raw: "unicum",
    pr: 2100,
    shipDamageRatio: 1.5,
    winRate: 60,
  },
  {
    raw: "super_unicum",
    pr: 2450,
    shipDamageRatio: 1.6,
    winRate: 65,
  },
] as const;

const COLORS: { [raw in RawRating]: ColorCode } = {
  bad: new ColorCode("#FE0E00"),
  below_avg: new ColorCode("#FE7903"),
  avg: new ColorCode("#FFC71F"),
  good: new ColorCode("#44B300"),
  very_good: new ColorCode("#318000"),
  great: new ColorCode("#02C9B3"),
  unicum: new ColorCode("#D042F3"),
  super_unicum: new ColorCode("#A00DC5"),
} as const;

const DISPLAY_NAMES: { [raw in RawRating]: string } = {
  bad: "Bad",
  below_avg: "Below Average",
  avg: "Average",
  good: "Good",
  very_good: "Very Good",
  great: "Great",
  unicum: "Unicum",
  super_unicum: "Super Unicum",
} as const;

export class Rating {
  constructor(public readonly raw: RawRating) {}

  static fromPR(value: number): Optional<Rating> {
    const satisfied = THRESHOLDS.findLast((it) => value >= it.pr);
    if (!satisfied) {
      return undefined;
    }
    return new Rating(satisfied.raw);
  }

  static fromWinRate(value: number): Optional<Rating> {
    const satisfied = THRESHOLDS.findLast((it) => value >= it.winRate);
    if (!satisfied) {
      return undefined;
    }
    return new Rating(satisfied.raw);
  }

  static fromShipDamage = (
    value: number,
    expected: number,
  ): Optional<Rating> => {
    if (expected <= 0) {
      return undefined;
    }
    const ratio = value / expected;

    const satisfied = THRESHOLDS.findLast((it) => ratio >= it.shipDamageRatio);
    if (!satisfied) {
      return undefined;
    }
    return new Rating(satisfied.raw);
  };

  static getThresholds(): RatingThreshold[] {
    return THRESHOLDS;
  }

  getColor(): ColorCode {
    return COLORS[this.raw];
  }

  getDisplayName(): string {
    return DISPLAY_NAMES[this.raw];
  }
}
