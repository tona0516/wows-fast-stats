import {
  isDigitKey,
  isOverallKey,
  isShipKey,
  isShipType,
  tierString,
} from "src/lib/util";

test("tierString - 正常系", () => {
  const values = [
    "I",
    "II",
    "III",
    "IV",
    "V",
    "VI",
    "VII",
    "VIII",
    "IX",
    "X",
    "★",
  ];

  for (let i = 0; i < values.length; i++) {
    const expected = values[i];
    expect(tierString(i + 1)).toBe(expected);
  }
});

test("tierString - 異常系", () => {
  expect(tierString(0)).toBe(undefined);
  expect(tierString(12)).toBe(undefined);
});

test("isShipType - 正常系", () => {
  const values = ["cv", "bb", "cl", "dd", "ss"];
  for (const value of values) {
    if (!isShipType(value)) fail();
  }
});

test("isShipType - 異常系", () => {
  const values = ["", "aux"];
  for (const value of values) {
    if (isShipType(value)) fail();
  }
});

test("isDigitKey - 正常系", () => {
  const values = ["battles", "hit_rate", "avg_tier"];
  for (const value of values) {
    if (!isDigitKey(value)) fail();
  }
});

test("isDigitKey - 異常系", () => {
  const values = ["", "invalid"];
  for (const value of values) {
    if (isDigitKey(value)) fail();
  }
});

test("isShipKey - 正常系", () => {
  const values = ["battles", "hit_rate"];
  for (const value of values) {
    if (!isShipKey(value)) fail();
  }
});

test("isShipKey - 異常系", () => {
  const values = ["", "avg_tier"];
  for (const value of values) {
    if (isShipKey(value)) fail();
  }
});

test("isOverallKey - 正常系", () => {
  const values = ["battles", "avg_tier"];
  for (const value of values) {
    if (!isOverallKey(value)) fail();
  }
});

test("isOverallKey - 異常系", () => {
  const values = ["", "hit_rate"];
  for (const value of values) {
    if (isOverallKey(value)) fail();
  }
});
