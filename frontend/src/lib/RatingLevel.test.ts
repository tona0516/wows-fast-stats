import { RatingInfo } from "src/lib/RatingLevel";
import { data } from "wailsjs/go/models";

test("ファクトリメソッド - 異常系", () => {
  const skillColorCode = new data.UCSkillColorCode();

  const pr = RatingInfo.fromPR(-1, skillColorCode);
  expect(pr).toBeUndefined();

  const damage1 = RatingInfo.fromDamage(16000, 0, skillColorCode);
  expect(damage1).toBeUndefined();

  const damage2 = RatingInfo.fromDamage(-1, 1, skillColorCode);
  expect(damage2).toBeUndefined();

  const winRate = RatingInfo.fromWinRate(-1, skillColorCode);
  expect(winRate).toBeUndefined();
});

test("ファクトリメソッド - 正常系", () => {
  const expectedTextColor = "#114514";

  const skillColorCode = new data.UCSkillColorCode({
    super_unicum: expectedTextColor,
  });

  const pr = RatingInfo.fromPR(2450, skillColorCode);
  expect(pr?.level).toBe("super_unicum");
  expect(pr?.textColorCode).toBe(expectedTextColor);

  const damage = RatingInfo.fromDamage(16000, 10000, skillColorCode);
  expect(damage?.level).toBe("super_unicum");
  expect(damage?.textColorCode).toBe(expectedTextColor);

  const winRate = RatingInfo.fromWinRate(100, skillColorCode);
  expect(winRate?.level).toBe("super_unicum");
  expect(winRate?.textColorCode).toBe(expectedTextColor);
});
