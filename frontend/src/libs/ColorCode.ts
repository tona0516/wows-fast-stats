import chroma from "chroma-js";
import { Theme } from "./Theme";

export class ColorCode {
  constructor(public readonly raw: string) {
    if (raw.match(/^#[A-Fa-f0-9]{6}$/) === null) {
      throw new Error(`Invalid color code: ${raw}`);
    }
  }

  getFixedTextColor(rate: number = 1.0): ColorCode {
    const chromaColor = chroma(this.raw);
    const fixed = Theme.isLight() ? chromaColor.darken(rate) : chromaColor;
    return new ColorCode(fixed.hex());
  }

  getFixedBgColor(rate: number = 1.0): ColorCode {
    const chromaColor = chroma(this.raw);
    const fixed = Theme.isLight() ? chromaColor.brighten(rate) : chromaColor;
    return new ColorCode(fixed.hex());
  }
}
