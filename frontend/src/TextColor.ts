import { LogDebug } from "../wailsjs/runtime/runtime";

namespace TextColor {
  export function pr(value: number): string {
    switch (true) {
      case value <= 0:
        return "";
      case value < 750:
        return "bad-text";
      case value < 1100:
        return "below-average-text";
      case value < 1350:
        return "average-text";
      case value < 1550:
        return "good-text";
      case value < 1750:
        return "very-good-text";
      case value < 2100:
        return "great-text";
      case value < 2450:
        return "unicum-text";
      case value >= 2450:
        return "super-unicum-text";
      default:
        return "";
    }
  }

  export function winRate(value: number): string {
    switch (true) {
      case value <= 0:
        return "";
      case value < 47:
        return "bad-text";
      case value < 50:
        return "below-average-text";
      case value < 52:
        return "average-text";
      case value < 54:
        return "good-text";
      case value < 56:
        return "very-good-text";
      case value < 60:
        return "great-text";
      case value < 65:
        return "unicum-text";
      case value >= 65:
        return "super-unicum-text";
      default:
        return "";
    }
  }

  export function shipType(value: string): string {
    switch (value) {
      case "AirCarrier":
        return "cv-text";
      case "Battleship":
        return "bb-text";
      case "Cruiser":
        return "cl-text";
      case "Destroyer":
        return "dd-text";
      case "Submarine":
        return "ss-text";
      default:
        return "";
    }
  }

  export function shipDamage(value: number, expected: number): string {
    const ratio = value / expected ?? 0;
    switch (true) {
      case ratio === 0:
        return "";
      case ratio < 0.6:
        return "bad-text";
      case ratio < 0.8:
        return "below-average-text";
      case ratio < 1.0:
        return "average-text";
      case ratio < 1.2:
        return "good-text";
      case ratio < 1.4:
        return "very-good-text";
      case ratio < 1.5:
        return "great-text";
      case ratio < 1.6:
        return "unicum-text";
      case ratio >= 1.6:
        return "super-unicum-text";
      default:
        return "";
    }
  }
}

export default TextColor;
