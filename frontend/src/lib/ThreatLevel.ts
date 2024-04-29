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

export class ThreatLevelInfo {
  constructor(
    readonly level: ThreatLevel,
    readonly textColorCode: string,
    readonly bgColorCode: string,
  ) {}

  static fromScore(score: number): ThreatLevelInfo | undefined {
    const tl = THREAT_LEVEL_DEFS.findLast((it) => score >= it.score);
    return tl?.info;
  }
}

const THREAT_LEVEL_COEF = 0.5;
export type ThreatLevelDef = {
  info: ThreatLevelInfo;
  score: number;
};
export const THREAT_LEVEL_DEFS: ThreatLevelDef[] = [
  {
    info: new ThreatLevelInfo(ThreatLevel.IR, "#FFFFFF", "#000000"),
    score: 0,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.R, "#FFFFFF", "#FF0000"),
    score: 8000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.O, "#331100", "#FFA500"),
    score: 13000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.Y, "#331100", "#FFFF00"),
    score: 19000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.G, "#FFFFFF", "#008000"),
    score: 25000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.B, "#FFFFFF", "#2255FF"),
    score: 32000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.I, "#FFFFFF", "#234794"),
    score: 35000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.V, "#FFFFFF", "#705DA8"),
    score: 40000 * THREAT_LEVEL_COEF,
  },
  {
    info: new ThreatLevelInfo(ThreatLevel.UV, "#000000", "#FFFFFF"),
    score: 44000 * THREAT_LEVEL_COEF,
  },
];
