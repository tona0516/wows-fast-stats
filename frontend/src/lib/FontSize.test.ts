import { FontSize } from "src/lib/FontSize";
import { data } from "wailsjs/go/models";

test("正常系", () => {
  const values = [
    { fontSize: "small", zoomRate: 0.85 },
    { fontSize: "medium", zoomRate: 1.0 },
    { fontSize: "large", zoomRate: 1.15 },
  ];

  for (const value of values) {
    expect(
      FontSize.getZoomRate(
        new data.UserConfigV2({ font_size: value.fontSize }),
      ),
    ).toBe(value.zoomRate);
  }
});

test("異常系", () => {
  const values = [
    { fontSize: "", zoomRate: 1.0 },
    { fontSize: "invalid", zoomRate: 1.0 },
  ];

  for (const value of values) {
    expect(
      FontSize.getZoomRate(
        new data.UserConfigV2({ font_size: value.fontSize }),
      ),
    ).toBe(value.zoomRate);
  }
});
