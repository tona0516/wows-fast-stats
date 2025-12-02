import chroma from "chroma-js";
import { Theme } from "./Theme";
import type { Rating, ShipType, TierGroup } from "./types";

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

export const TIER_GROUP_COLORS: { [tierGroup in TierGroup]: ColorCode } = {
  low: new ColorCode("#187FC4"),
  middle: new ColorCode("#AACF52"),
  high: new ColorCode("#EA5532"),
} as const;

export const SHIP_TYPE_COLORS: { [type in ShipType]: ColorCode } = {
  ss: new ColorCode("#233B8B"),
  dd: new ColorCode("#D9760F"),
  cl: new ColorCode("#27853F"),
  bb: new ColorCode("#CA1028"),
  cv: new ColorCode("#5E2883"),
} as const;

export const RATING_COLORS: { [rating in Rating]: ColorCode } = {
  bad: new ColorCode("#FE0E00"),
  below_avg: new ColorCode("#FE7903"),
  avg: new ColorCode("#FFC71F"),
  good: new ColorCode("#44B300"),
  very_good: new ColorCode("#318000"),
  great: new ColorCode("#02C9B3"),
  unicum: new ColorCode("#D042F3"),
  super_unicum: new ColorCode("#A00DC5"),
} as const;

type ColorPair = {
  readonly text: ColorCode;
  readonly background: ColorCode;
};

export const THREAT_LEVEL_COLORS: { [rank: string]: ColorPair } = {
  ir: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#000000") },
  r: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#FF0000") },
  o: { text: new ColorCode("#331100"), background: new ColorCode("#FFA500") },
  y: { text: new ColorCode("#331100"), background: new ColorCode("#FFFF00") },
  g: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#008000") },
  b: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#2255FF") },
  i: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#234794") },
  v: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#705DA8") },
  uv: { text: new ColorCode("#FFFFFF"), background: new ColorCode("#800080") },
} as const;
