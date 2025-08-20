import chroma from "chroma-js";
import { Theme } from "./Theme";

export class ColorCode {
  static readonly PLAYER_BG_FIDED_RATE = 2.25;
  static readonly SHIP_TYPE_BG_FIXED_RATE = 2;

  constructor(public readonly raw: string) {
    if (raw.match(/^#[A-Fa-f0-9]{6}$/) === null) {
      throw new Error(`Invalid color code: ${raw}`);
    }
  }

  getFixedTextColor(rate: number = 1.0): ColorCode {
    const chromaColor = chroma(this.raw);
    const fixed = Theme.isLighter()
      ? chromaColor.darken(rate)
      : chromaColor.brighten(rate);
    return new ColorCode(fixed.hex());
  }

  getFixedBgColor(rate: number = 1.0): ColorCode {
    const chromaColor = chroma(this.raw);
    const fixed = Theme.isLighter()
      ? chromaColor.brighten(rate)
      : chromaColor.darken(rate);
    return new ColorCode(fixed.hex());
  }
}
