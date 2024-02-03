import { Rating } from "src/lib/Rating";
import { model } from "wailsjs/go/models";

test("none", () => {
  const expectedTextColor = "";

  const skillColorCode = new model.UCSkillColorCode();

  // pattern 1: rating is undefined
  const converter = new Rating(undefined, skillColorCode);
  expect(converter.colorCode()).toBe(expectedTextColor);

  // pattern 2: expected is zero
  const damage = Rating.fromDamage(16000, 0, skillColorCode);
  expect(damage.level).toBeUndefined();
  expect(damage.colorCode()).toBe(expectedTextColor);
});

test("factory", () => {
  const expectedTextColor = "#114514";

  const skillColorCode = new model.UCSkillColorCode({
    super_unicum: expectedTextColor,
  });

  const pr = Rating.fromPR(2450, skillColorCode);
  expect(pr.level).toBe("super_unicum");
  expect(pr.colorCode()).toBe(expectedTextColor);

  const damage = Rating.fromDamage(16000, 10000, skillColorCode);
  expect(damage.level).toBe("super_unicum");
  expect(damage.colorCode()).toBe(expectedTextColor);

  const winRate = Rating.fromWinRate(100, skillColorCode);
  expect(winRate.level).toBe("super_unicum");
  expect(winRate.colorCode()).toBe(expectedTextColor);
});
