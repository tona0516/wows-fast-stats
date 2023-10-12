import { isOverallKey, isShipKey, isShipType, tierString } from "src/lib/util";

test("tierString", () => {
  const expecteds = [
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
  ];
  expecteds.forEach((expected, i) => {
    expect(tierString(i)).toBe(expected);
  });
});

test("isShipType", () => {
  const validTypes = ["cv", "bb", "cl", "dd", "ss"];
  validTypes.forEach((type) => {
    if (!isShipType(type)) {
      fail();
    }
  });

  const invalidTypes = ["", "aux"];
  invalidTypes.forEach((type) => {
    if (isShipType(type)) {
      fail();
    }
  });
});

test("isShipKey", () => {
  const validKeys = ["battles", "hit_rate"];
  validKeys.forEach((key) => {
    if (!isShipKey(key)) {
      fail();
    }
  });

  if (isShipKey("avg_tier")) {
    fail();
  }
});

test("isOverallKey", () => {
  const validKeys = ["battles", "avg_tier"];
  validKeys.forEach((key) => {
    if (!isOverallKey(key)) {
      fail();
    }
  });

  if (isOverallKey("hit_rate")) {
    fail();
  }
});
