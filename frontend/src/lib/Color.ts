import chroma from "chroma-js";
import { AppFunctions } from "./AppFunctions";
import type { RatingLevel } from "./RatingLevel";
import type { ThreatLevel } from "./ThreatLevel";
import type { ShipType } from "./types";

export interface ColorPair {
  text: string;
  background: string;
}

export namespace Color {
  export namespace ShipType {
    export const getDefault = (type: ShipType): string => {
      switch (type) {
        case "cv":
          return "#5E2883";
        case "bb":
          return "#CA1028";
        case "cl":
          return "#27853F";
        case "dd":
          return "#D9760F";
        case "ss":
          return "#233B8B";
      }
    };

    export const getFixed = (type: ShipType): ColorPair => {
      const colorCode = getDefault(type);
      return getFixedColorPair(colorCode);
    };
  }

  export namespace ThreatLevel {
    export const getDefault = (level: ThreatLevel): ColorPair => {
      switch (level) {
        case "ir":
          return { text: "#FFFFFF", background: "#000000" };
        case "r":
          return { text: "#FFFFFF", background: "#FF0000" };
        case "o":
          return { text: "#331100", background: "#FFA500" };
        case "y":
          return { text: "#331100", background: "#FFFF00" };
        case "g":
          return { text: "#FFFFFF", background: "#008000" };
        case "b":
          return { text: "#FFFFFF", background: "#2255FF" };
        case "i":
          return { text: "#FFFFFF", background: "#234794" };
        case "v":
          return { text: "#FFFFFF", background: "#705DA8" };
        case "uv":
          return { text: "#FFFFFF", background: "#800080" };
      }
    };
  }

  export namespace Rating {
    export const getDefault = (level: RatingLevel): string => {
      switch (level) {
        case "bad":
          return "#FE0E00";
        case "below_avg":
          return "#FE7903";
        case "avg":
          return "#FFC71F";
        case "good":
          return "#44B300";
        case "very_good":
          return "#318000";
        case "great":
          return "#02C9B3";
        case "unicum":
          return "#D042F3";
        case "super_unicum":
          return "#A00DC5";
      }
    };

    export const getFixed = (level: RatingLevel): ColorPair => {
      const colorCode = getDefault(level);
      return getFixedColorPair(colorCode);
    };
  }
}

const getFixedColorPair = (colorCode: string): ColorPair => {
  const chromaColor = chroma(colorCode);
  const colorRate = 1.5;
  const brighten = chromaColor.brighten(colorRate).hex();
  const darken = chromaColor.darken(colorRate).hex();

  if (AppFunctions.isLighter()) {
    return { text: darken, background: brighten };
  }
  return { text: brighten, background: darken };
};
