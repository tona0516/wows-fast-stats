import { RatingfGenerator } from "./RatingLevel";

test("ファクトリメソッド - 異常系", () => {
  const pr = RatingfGenerator.fromPR(-1);
  expect(pr).toBeUndefined();

  const damage1 = RatingfGenerator.fromDamage(16000, 0);
  expect(damage1).toBeUndefined();

  const damage2 = RatingfGenerator.fromDamage(-1, 1);
  expect(damage2).toBeUndefined();

  const winRate = RatingfGenerator.fromWinRate(-1);
  expect(winRate).toBeUndefined();
});

test("ファクトリメソッド - 正常系", () => {
  const pr = RatingfGenerator.fromPR(2450);
  expect(pr?.level).toBe("super_unicum");
  expect(pr?.color).toBeDefined();

  const damage = RatingfGenerator.fromDamage(16000, 10000);
  expect(damage?.level).toBe("super_unicum");
  expect(damage?.color).toBeDefined();

  const winRate = RatingfGenerator.fromWinRate(100);
  expect(winRate?.level).toBe("super_unicum");
  expect(damage?.color).toBeDefined();
});
