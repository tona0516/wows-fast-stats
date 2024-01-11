import { FontSize } from "src/lib/FontSize";
import { model } from "wailsjs/go/models";

test("正常系", () => {
  [
    { fontSize: "small", zoomRate: 0.85 },
    { fontSize: "medium", zoomRate: 1.0 },
    { fontSize: "large", zoomRate: 1.15 },
  ].forEach((it) => {
    expect(
      FontSize.getZoomRate(new model.UserConfig({ font_size: it.fontSize })),
    ).toBe(it.zoomRate);
  });
});

test("異常系", () => {
  [
    { fontSize: "", zoomRate: 1.0 },
    { fontSize: "invalid", zoomRate: 1.0 },
  ].forEach((it) => {
    expect(
      FontSize.getZoomRate(new model.UserConfig({ font_size: it.fontSize })),
    ).toBe(it.zoomRate);
  });
});
