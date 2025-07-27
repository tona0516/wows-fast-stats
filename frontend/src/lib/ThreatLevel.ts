import { Color, type ColorPair } from "./Color";

export type ThreatLevel = "ir" | "r" | "o" | "y" | "g" | "b" | "i" | "v" | "uv";

export class ThreatLevelGenerator {
  private constructor() {}

  static fromScore(score: number): ThreatLevelInfo | undefined {
    return THREAT_LEVEL_DEFS.findLast((it) => score >= it.score);
  }
}
export interface ThreatLevelInfo {
  level: ThreatLevel;
  score: number;
  color: ColorPair;
}

const THREAT_LEVEL_COEF = 0.5;
export const THREAT_LEVEL_DEFS: ThreatLevelInfo[] = [
  {
    level: "ir",
    score: 0,
    color: Color.ThreatLevel.getDefault("ir"),
  },
  {
    level: "r",
    score: 8000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("r"),
  },
  {
    level: "o",
    score: 13000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("o"),
  },
  {
    level: "y",
    score: 19000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("y"),
  },
  {
    level: "g",
    score: 25000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("g"),
  },
  {
    level: "b",
    score: 32000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("b"),
  },
  {
    level: "i",
    score: 35000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("i"),
  },
  {
    level: "v",
    score: 40000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("v"),
  },
  {
    level: "uv",
    score: 44000 * THREAT_LEVEL_COEF,
    color: Color.ThreatLevel.getDefault("uv"),
  },
];
