import { RATING_FACTORS } from "src/lib/rating/RatingConst";
import type { Rating } from "src/lib/types";
import type { domain } from "wailsjs/go/models";

const NONE_COLOR = "#000000";

export class RatingConverter {
  constructor(
    public rating: Rating | undefined,
    private config: domain.UserConfig,
  ) {}

  getTextColorCode(): string {
    if (!this.rating) {
      return "";
    }
    return this.config.custom_color.skill.text[this.rating];
  }

  getBgColorCode(): string {
    if (!this.rating) {
      return NONE_COLOR;
    }
    return this.config.custom_color.skill.background[this.rating];
  }
}

export namespace RatingConverterFactory {
  export const fromPR = (
    value: number,
    config: domain.UserConfig,
  ): RatingConverter => {
    const rf = RATING_FACTORS.findLast(
      (it) => value >= 0 && value >= it.pr.min,
    );
    return new RatingConverter(rf?.level, config);
  };

  export const fromDamage = (
    value: number,
    expected: number,
    config: domain.UserConfig,
  ): RatingConverter => {
    if (expected === 0) {
      return new RatingConverter(undefined, config);
    }

    const ratio = value / expected;
    const rf = RATING_FACTORS.findLast((it) => ratio >= it.damage.min);
    return new RatingConverter(rf?.level, config);
  };

  export const fromWinRate = (
    value: number,
    config: domain.UserConfig,
  ): RatingConverter => {
    const rf = RATING_FACTORS.findLast((it) => value >= it.winRate.min);
    return new RatingConverter(rf?.level, config);
  };
}
