import type { domain } from "../wailsjs/go/models";
import { SkillLevel } from "./enums";

const skillLevels = Object.values(SkillLevel);

interface Range {
  skillLevel: SkillLevel;
  max: number;
}

class Ranges {
  values: Array<Range>;

  constructor(maxs: number[]) {
    if (maxs.length !== skillLevels.length) {
      throw Error("Lengths of maxs and skillLevels do not match");
    }

    this.values = new Array(skillLevels.length);
    for (let i = 0; i < this.values.length; i++) {
      this.values[i] = { skillLevel: skillLevels[i], max: maxs[i] };
    }
  }
}

const prRange = new Ranges([
  750,
  1100,
  1350,
  1550,
  1750,
  2100,
  2450,
  Number.MAX_VALUE,
]);

const damageRatioRange = new Ranges([
  0.6,
  0.8,
  1.0,
  1.2,
  1.4,
  1.5,
  1.6,
  Number.MAX_VALUE,
]);

const winRateRange = new Ranges([47, 50, 52, 54, 56, 60, 65, Number.MAX_VALUE]);

export class SkillLevelConverter {
  skillColor: domain.SkillColor;
  skillLevel?: SkillLevel;

  private constructor(skillColor: domain.SkillColor, skillLevel?: SkillLevel) {
    this.skillColor = skillColor;
    this.skillLevel = skillLevel;
  }

  static fromPR(
    value: number,
    skillColor: domain.SkillColor
  ): SkillLevelConverter {
    const range = prRange.values.find((it) => value >= 0 && value < it.max);
    return range
      ? new SkillLevelConverter(skillColor, range.skillLevel)
      : new SkillLevelConverter(skillColor);
  }

  static fromDamage(
    value: number,
    expected: number,
    skillColor: domain.SkillColor
  ): SkillLevelConverter {
    const ratio = value / expected ?? 0;
    const range = damageRatioRange.values.find((it) => ratio < it.max);
    return range
      ? new SkillLevelConverter(skillColor, range.skillLevel)
      : new SkillLevelConverter(skillColor);
  }

  static fromWinRate(
    value: number,
    skillColor: domain.SkillColor
  ): SkillLevelConverter {
    const range = winRateRange.values.find((it) => value < it.max);
    return range
      ? new SkillLevelConverter(skillColor, range.skillLevel)
      : new SkillLevelConverter(skillColor);
  }

  toTextColorCode(): string {
    return this.skillColor.text[this.skillLevel] ?? "";
  }

  toBgColorCode(): string {
    return this.skillColor.background[this.skillLevel] ?? "";
  }
}
