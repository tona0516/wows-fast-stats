export enum ThreatLevel {
  IR = "IR",
  R = "R",
  O = "O",
  Y = "Y",
  G = "G",
  B = "B",
  I = "I",
  V = "V",
  UV = "UV",
}

export class ThreatLevelGenerator {
  private constructor() {}

  static fromScore(score: number): ThreatLevelInfo | undefined {
    return THREAT_LEVEL_DEFS.findLast((it) => score >= it.score);
  }
}
export interface ThreatLevelInfo {
  level: ThreatLevel;
  score: number;
  textColor: string;
  bgColor: string;
}

const THREAT_LEVEL_COEF = 0.5;
export const THREAT_LEVEL_DEFS: ThreatLevelInfo[] = [
  {
    level: ThreatLevel.IR,
    textColor: "#FFFFFF",
    bgColor: "#000000",
    score: 0,
  },
  {
    level: ThreatLevel.R,
    textColor: "#FFFFFF",
    bgColor: "#FF0000",
    score: 8000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.O,
    textColor: "#331100",
    bgColor: "#FFA500",
    score: 13000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.Y,
    textColor: "#331100",
    bgColor: "#FFFF00",
    score: 19000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.G,
    textColor: "#FFFFFF",
    bgColor: "#2255FF",
    score: 25000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.B,
    textColor: "#FFFFFF",
    bgColor: "#FFA500",
    score: 32000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.I,
    textColor: "#FFFFFF",
    bgColor: "#234794",
    score: 35000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.V,
    textColor: "#FFFFFF",
    bgColor: "#705DA8",
    score: 40000 * THREAT_LEVEL_COEF,
  },
  {
    level: ThreatLevel.UV,
    textColor: "#000000",
    bgColor: "#FFFFFF",
    score: 44000 * THREAT_LEVEL_COEF,
  },
];
