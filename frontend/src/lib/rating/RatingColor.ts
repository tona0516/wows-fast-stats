import type { Rating } from "src/lib/types";
import type { model } from "wailsjs/go/models";

export class RatingAdapter {
  constructor(
    public rating: Rating | undefined,
    private config: model.UserConfig,
  ) {}

  getTextColorCode(): string {
    if (!this.rating) {
      return "";
    }
    return this.config.color.skill.text[this.rating];
  }
}
