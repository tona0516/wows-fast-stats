import { RatingAdapter } from "src/lib/rating/RatingColor";
import { RatingColorFactory } from "src/lib/rating/RatingColorFactory";
import { model } from "wailsjs/go/models";

test("none", () => {
  const expectedTextColor = "";

  const config = new model.UserConfig();

  // pattern 1: rating is undefined
  const converter = new RatingAdapter(undefined, config);
  expect(converter.getTextColorCode()).toBe(expectedTextColor);

  // pattern 2: expected is zero
  const damage = RatingColorFactory.fromDamage(16000, 0, config);
  expect(damage.rating).toBeUndefined();
  expect(damage.getTextColorCode()).toBe(expectedTextColor);
});

test("factory", () => {
  const expectedTextColor = "#000001";
  const expectedBgColor = "#000002";

  const config = new model.UserConfig({
    custom_color: {
      skill: {
        text: { super_unicum: expectedTextColor },
        background: { super_unicum: expectedBgColor },
      },
    },
  });

  const pr = RatingColorFactory.fromPR(2450, config);
  expect(pr.rating).toBe("super_unicum");
  expect(pr.getTextColorCode()).toBe(expectedTextColor);

  const damage = RatingColorFactory.fromDamage(16000, 10000, config);
  expect(damage.rating).toBe("super_unicum");
  expect(damage.getTextColorCode()).toBe(expectedTextColor);

  const winRate = RatingColorFactory.fromWinRate(100, config);
  expect(winRate.rating).toBe("super_unicum");
  expect(winRate.getTextColorCode()).toBe(expectedTextColor);
});
