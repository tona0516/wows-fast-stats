import chroma from "chroma-js";
import type { RatingLevel } from "./RatingLevel";
import { Theme } from "./Theme";
import type { ThreatLevel } from "./ThreatLevel";
import type { ShipType } from "./types";

const COLOR_RATE = 1.5

export interface ColorPair {
  text: string;
  background: string;
}

export namespace Color {
  export namespace ShipType {
    const TYPES: Map<ShipType, string> = new Map([
      ["cv", "#5E2883"],
      ["bb", "#CA1028"],
      ["cl", "#27853F"],
      ["dd", "#D9760F"],
      ["ss", "#233B8B"],
    ]);

    export const getDefault = (type: ShipType): string | undefined => {
      return TYPES.get(type);
    };

    export const getFixed = (type: ShipType): ColorPair | undefined => {
      const colorCode = TYPES.get(type);
      if (!colorCode) {
        return undefined;
      }

      const chromaColor = chroma(colorCode);
      const brighten = chromaColor.brighten(COLOR_RATE).hex();
      const darken = chromaColor.darken(COLOR_RATE).hex();

      if (Theme.isLighter()) {
        return { text: darken, background: brighten };
      }
      return { text: brighten, background: darken };
    };
  }

  export namespace ThreatLevel {
    const TYPES: Map<ThreatLevel, ColorPair> = new Map([
      ["ir", { text: "#FFFFFF", background: "#000000" }],
      ["r", { text: "#FFFFFF", background: "#FF0000" }],
      ["o", { text: "#331100", background: "#FFA500" }],
      ["y", { text: "#331100", background: "#FFFF00" }],
      ["g", { text: "#FFFFFF", background: "#008000" }],
      ["b", { text: "#FFFFFF", background: "#2255FF" }],
      ["i", { text: "#FFFFFF", background: "#234794" }],
      ["v", { text: "#FFFFFF", background: "#705DA8" }],
      ["uv", { text: "#FFFFFF", background: "#800080" }],
    ]);

    export const getDefault = (level: ThreatLevel): ColorPair | undefined => {
      return TYPES.get(level);
    };
  }

  export namespace Rating {
    const TYPES: Map<RatingLevel, string> = new Map([
      ["bad", "#FE0E00"],
      ["below_avg", "#FE7903"],
      ["avg", "#FFC71F"],
      ["good", "#44B300"],
      ["very_good", "#318000"],
      ["great", "#02C9B3"],
      ["unicum", "#D042F3"],
      ["super_unicum", "#A00DC5"],
    ]);

    export const getDefault = (level: RatingLevel): string | undefined => {
      return TYPES.get(level);
    };

    export const getFixed = (level: RatingLevel): ColorPair | undefined => {
      const colorCode = TYPES.get(level);
      if (!colorCode) {
        return undefined;
      }

      const chromaColor = chroma(colorCode);

      const brighten = chromaColor.brighten(COLOR_RATE).hex();
      const darken = chromaColor.darken(COLOR_RATE).hex();

      if (Theme.isLighter()) {
        return { text: darken, background: brighten };
      }
      return { text: brighten, background: darken };
    };
  }
}
