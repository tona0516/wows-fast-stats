import { Const, type SkillLevelItem } from "src/Const";
import type { domain } from "wailsjs/go/models";

export class SkillLevelConverter {
  skillColor: domain.SkillColor;
  skillLevelItem?: SkillLevelItem;

  private constructor(
    skillColor: domain.SkillColor,
    skillLevelItem?: SkillLevelItem
  ) {
    this.skillColor = skillColor;
    this.skillLevelItem = skillLevelItem;
  }

  static fromPR(value: number, color: domain.SkillColor): SkillLevelConverter {
    const level = Const.SKILL_LEVELS.find(
      (it) => value >= 0 && value < it.maxPR
    );
    return level
      ? new SkillLevelConverter(color, level)
      : new SkillLevelConverter(color);
  }

  static fromDamage(
    value: number,
    expected: number,
    color: domain.SkillColor
  ): SkillLevelConverter {
    const ratio = value / expected ?? 0;
    const level = Const.SKILL_LEVELS.find((it) => ratio < it.maxDamage);
    return level
      ? new SkillLevelConverter(color, level)
      : new SkillLevelConverter(color);
  }

  static fromWinRate(
    value: number,
    color: domain.SkillColor
  ): SkillLevelConverter {
    const level = Const.SKILL_LEVELS.find((it) => value < it.maxWin);
    return level
      ? new SkillLevelConverter(color, level)
      : new SkillLevelConverter(color);
  }

  toTextColorCode(): string {
    if (!this.skillLevelItem) {
      return "";
    }
    return this.skillColor.text[this.skillLevelItem.level];
  }

  toBgColorCode(): string {
    if (!this.skillLevelItem) {
      return "";
    }
    return this.skillColor.background[this.skillLevelItem.level];
  }
}
