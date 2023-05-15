import Const from "./Const";

type Rank =
  | ""
  | "bad"
  | "belowAvg"
  | "avg"
  | "good"
  | "veryGood"
  | "great"
  | "unicum"
  | "superUnicum";

type Range = {
  rank: Rank;
  max: number;
};

export default class RankConverter {
  rank: Rank;

  private constructor(rank: Rank) {
    this.rank = rank;
  }

  static fromPR(value: number): RankConverter {
    const ranges: Range[] = [
      { rank: "bad", max: 750 },
      { rank: "belowAvg", max: 1100 },
      { rank: "avg", max: 1350 },
      { rank: "good", max: 1550 },
      { rank: "veryGood", max: 1750 },
      { rank: "great", max: 2100 },
      { rank: "unicum", max: 2450 },
      { rank: "superUnicum", max: Number.MAX_VALUE },
    ];

    for (const range of ranges) {
      if (value > 0 && value < range.max) {
        return new RankConverter(range.rank);
      }
    }

    return new RankConverter("");
  }

  static fromDamage(value: number, expected: number): RankConverter {
    const ratio = value / expected ?? 0;

    const ranges: Range[] = [
      { rank: "bad", max: 0.6 },
      { rank: "belowAvg", max: 0.8 },
      { rank: "avg", max: 1.0 },
      { rank: "good", max: 1.2 },
      { rank: "veryGood", max: 1.4 },
      { rank: "great", max: 1.5 },
      { rank: "unicum", max: 1.6 },
      { rank: "superUnicum", max: Number.MAX_VALUE },
    ];

    for (const range of ranges) {
      if (ratio < range.max) {
        return new RankConverter(range.rank);
      }
    }

    return new RankConverter("");
  }

  static fromWinRate(value: number): RankConverter {
    const ranges: Range[] = [
      { rank: "bad", max: 47 },
      { rank: "belowAvg", max: 50 },
      { rank: "avg", max: 52 },
      { rank: "good", max: 54 },
      { rank: "veryGood", max: 56 },
      { rank: "great", max: 60 },
      { rank: "unicum", max: 65 },
      { rank: "superUnicum", max: Number.MAX_VALUE },
    ];

    for (const range of ranges) {
      if (value < range.max) {
        return new RankConverter(range.rank);
      }
    }

    return new RankConverter("");
  }

  toTextColorCode(): string {
    return Const.RANK_TEXT_COLORS[this.rank];
  }

  toBgColorCode(): string {
    return Const.RANK_BG_COLORS[this.rank];
  }
}
