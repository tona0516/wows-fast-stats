import {
  isDigitKey,
  isOverallKey,
  isShipKey,
  isShipType,
  tierString,
} from "src/lib/util";

test("tierString", () => {
  [
    "",
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
  ].forEach((expected, i) => {
    expect(tierString(i)).toBe(expected);
  });
});

test("isShipType - 正常系", () => {
  ["cv", "bb", "cl", "dd", "ss"].forEach((type) => {
    if (!isShipType(type)) fail();
  });
});

test("isShipType - 異常系", () => {
  ["", "aux"].forEach((type) => {
    if (isShipType(type)) fail();
  });
});

test("isDigitKey - 正常系", () => {
  ["battles", "hit_rate", "avg_tier"].forEach((key) => {
    if (!isDigitKey(key)) fail();
  });
});

test("isDigitKey - 異常系", () => {
  ["", "invalid"].forEach((key) => {
    if (isDigitKey(key)) fail();
  });
});

test("isShipKey - 正常系", () => {
  ["battles", "hit_rate"].forEach((key) => {
    if (!isShipKey(key)) fail();
  });
});

test("isShipKey - 異常系", () => {
  ["", "avg_tier"].forEach((key) => {
    if (isShipKey(key)) fail();
  });
});

test("isOverallKey - 正常系", () => {
  ["battles", "avg_tier"].forEach((key) => {
    if (!isOverallKey(key)) fail();
  });
});

test("isOverallKey - 異常系", () => {
  ["", "hit_rate"].forEach((key) => {
    if (isOverallKey(key)) fail();
  });
});
