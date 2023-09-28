import type { Rating } from "src/lib/types";

export class RatingFactor {
  level: Rating;
  tier: number;
  shipType: string;
  minPR: number;
  maxPR: number;
  minDamage: number;
  maxDamage: number;
  minWin: number;
  maxWin: number;

  constructor(
    level: Rating,
    tier: number,
    shipType: string,
    minPR: number,
    maxPR: number,
    minDamage: number,
    maxDamage: number,
    minWin: number,
    maxWin: number,
  ) {
    this.level = level;
    this.tier = tier;
    this.shipType = shipType;
    this.minPR = minPR;
    this.maxPR = maxPR;
    this.minDamage = minDamage;
    this.maxDamage = maxDamage;
    this.minWin = minWin;
    this.maxWin = maxWin;
  }
}
