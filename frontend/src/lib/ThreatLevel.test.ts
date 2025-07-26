import {
  THREAT_LEVEL_DEFS,
  ThreatLevel,
  ThreatLevelGenerator,
} from "src/lib/ThreatLevel";

test("fromScore - 異常系", () => {
  expect(ThreatLevelGenerator.fromScore(-1)).toBeUndefined();
});

test("fromScore - 正常系", () => {
  const instance = ThreatLevelGenerator.fromScore(9500);
  expect(instance?.level).toBe(ThreatLevel.Y);
  expect(instance?.score).toBeDefined();
  expect(instance?.textColor).toBeDefined();
  expect(instance?.bgColor).toBeDefined();
});

test("THREAT_LEVEL_DEFSがスコアの昇順で定義されている", () => {
  const sorted = THREAT_LEVEL_DEFS.toSorted((a, b) => a.score - b.score);
  for (let i = 0; i < THREAT_LEVEL_DEFS.length; i++) {
    expect(THREAT_LEVEL_DEFS[i].score).toBe(sorted[i].score);
  }
});
