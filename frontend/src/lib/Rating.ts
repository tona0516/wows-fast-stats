import { RATING_DEFS } from "src/lib/RatingDef";
import type { RatingLevel } from "src/lib/types";
import { model } from "wailsjs/go/models";

export class Rating {
  constructor(
    public level: RatingLevel | undefined,
    private codeColor: model.UCSkillColorCode,
  ) {}

  colorCode(): string {
    return this.level ? this.codeColor[this.level] : "";
  }

  static fromPR(value: number, colorCode: model.UCSkillColorCode): Rating {
    const rf = RATING_DEFS.findLast((it) => value >= it.pr.min);
    return new Rating(rf?.level, colorCode);
  }

  static fromWinRate(value: number, colorCode: model.UCSkillColorCode): Rating {
    const rf = RATING_DEFS.findLast((it) => value >= it.winRate.min);
    return new Rating(rf?.level, colorCode);
  }

  static fromDamage(
    value: number,
    expected: number,
    colorCode: model.UCSkillColorCode,
  ): Rating {
    if (expected === 0) {
      return new Rating(undefined, colorCode);
    }

    const ratio = value / expected;
    const rf = RATING_DEFS.findLast((it) => ratio >= it.damage.min);
    return new Rating(rf?.level, colorCode);
  }
}
