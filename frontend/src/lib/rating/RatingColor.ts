import type { Rating } from "src/lib/types";
import type { model } from "wailsjs/go/models";

const NONE_COLOR = "#000000";

export class RatingAdapter {
  constructor(
    public rating: Rating | undefined,
    private config: model.UserConfig,
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
