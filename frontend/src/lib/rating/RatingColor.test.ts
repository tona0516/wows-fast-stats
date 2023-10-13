import { RatingAdapter } from "src/lib/rating/RatingColor";
import { RatingColorFactory } from "src/lib/rating/RatingColorFactory";
import { domain } from "wailsjs/go/models";

test("none", () => {
  const expectedTextColor = "";
  const expectedBgColor = "#000000";

  const config = new domain.UserConfig();

  // pattern 1: rating is undefined
  const converter = new RatingAdapter(undefined, config);
  expect(converter.getTextColorCode()).toBe(expectedTextColor);
  expect(converter.getBgColorCode()).toBe(expectedBgColor);

  // pattern 2: expected is zero
  const damage = RatingColorFactory.fromDamage(16000, 0, config);
  expect(damage.rating).toBeUndefined();
  expect(damage.getTextColorCode()).toBe(expectedTextColor);
  expect(damage.getBgColorCode()).toBe(expectedBgColor);
});

test("factory", () => {
  const expectedTextColor = "#000001";
  const expectedBgColor = "#000002";

  const config = new domain.UserConfig({
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
  expect(pr.getBgColorCode()).toBe(expectedBgColor);

  const damage = RatingColorFactory.fromDamage(16000, 10000, config);
  expect(damage.rating).toBe("super_unicum");
  expect(damage.getTextColorCode()).toBe(expectedTextColor);
  expect(damage.getBgColorCode()).toBe(expectedBgColor);

  const winRate = RatingColorFactory.fromWinRate(100, config);
  expect(winRate.rating).toBe("super_unicum");
  expect(winRate.getTextColorCode()).toBe(expectedTextColor);
  expect(winRate.getBgColorCode()).toBe(expectedBgColor);
});
