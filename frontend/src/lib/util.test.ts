import { isShipType, tierString } from "src/lib/util";

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
  expect(tierString(0)).toBe("");
  expect(tierString(12)).toBe("");
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
