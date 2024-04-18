import type { data } from "wailsjs/go/models";

export enum FontSize {
  XX_SMALL = "xx-small",
  X_SMALL = "x-small",
  SMALL = "small",
  MEDIUM = "medium",
  LARGE = "large",
  X_LARGE = "x-large",
  XX_LARGE = "xx-large",
}

const ZOOM_RATIO = new Map<FontSize, number>([
  [FontSize.XX_SMALL, 0.55],
  [FontSize.X_SMALL, 0.7],
  [FontSize.SMALL, 0.85],
  [FontSize.MEDIUM, 1.0],
  [FontSize.LARGE, 1.15],
  [FontSize.X_LARGE, 1.3],
  [FontSize.XX_LARGE, 1.55],
]);

export namespace FontSize {
  export const getZoomRate = (config: data.UserConfigV2): number => {
    const fontSize = config.font_size as FontSize;
    const zoomRatio = ZOOM_RATIO.get(fontSize);

    return zoomRatio || 1.0;
  };
}
