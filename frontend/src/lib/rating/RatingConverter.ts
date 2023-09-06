import { ratingFactors } from "src/lib/rating/RatingConst";
import type { Rating } from "src/lib/types";
import type { domain } from "wailsjs/go/models";

export class RatingConverter {
  constructor(
    private rating: Rating | undefined,
    private config: domain.UserConfig,
  ) {}

  textColorCode(): string {
    if (!this.rating) {
      return "";
    }
    return this.config.custom_color.skill.text[this.rating];
  }

  bgColorCode(): string {
    if (!this.rating) {
      return "";
    }
    return this.config.custom_color.skill.background[this.rating];
  }
}

export class RatingConverterFactory {
  static fromPR(value: number, config: domain.UserConfig): RatingConverter {
    const item = ratingFactors().findLast(
      (it) => value >= 0 && value >= it.minPR,
    );
    return new RatingConverter(item?.level, config);
  }

  static fromDamage(
    value: number,
    expected: number,
    config: domain.UserConfig,
  ): RatingConverter {
    const ratio = value / expected ?? 0;
    const item = ratingFactors().findLast((it) => ratio >= it.minDamage);
    return new RatingConverter(item?.level, config);
  }

  static fromWinRate(
    value: number,
    config: domain.UserConfig,
  ): RatingConverter {
    const item = ratingFactors().findLast((it) => value >= it.minWin);
    return new RatingConverter(item?.level, config);
  }
}
