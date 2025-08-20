import { ColorCode } from "./ColorCode";
import type { Optional } from "./types";

type RawThreatLevel = "ir" | "r" | "o" | "y" | "g" | "b" | "i" | "v" | "uv";

type ThreatLevelThreshold = {
  readonly level: RawThreatLevel;
  readonly threshold: number;
};

const COEF = 0.5;
const THRESHOLDS: ThreatLevelThreshold[] = [
  { level: "ir", threshold: 0 },
  { level: "r", threshold: 8000 * COEF },
  { level: "o", threshold: 13000 * COEF },
  { level: "y", threshold: 19000 * COEF },
  { level: "g", threshold: 25000 * COEF },
  { level: "b", threshold: 32000 * COEF },
  { level: "i", threshold: 35000 * COEF },
  { level: "v", threshold: 40000 * COEF },
  { level: "uv", threshold: 44000 * COEF },
] as const;

type ColorPair = {
  readonly text: ColorCode;
  readonly background: ColorCode;
};

const COLORS: { [raw in RawThreatLevel]: ColorPair } = {
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

export class ThreatLevel {
  constructor(public readonly raw: RawThreatLevel) {}

  static fromScore(score: number): Optional<ThreatLevel> {
    const satisfied = THRESHOLDS.findLast((it) => score >= it.threshold);
    if (!satisfied) {
      return undefined;
    }

    return new ThreatLevel(satisfied.level);
  }

  getColor(): ColorPair {
    return COLORS[this.raw];
  }
}
