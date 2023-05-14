import Const from "./Const";

namespace TextColor {
  export function prBG(value: number): string {
    switch (true) {
      case value <= 0:
        return "";
      case value < 750:
        return Const.PR_BG_COLORS.bad;
      case value < 1100:
        return Const.PR_BG_COLORS.belowAvg;
      case value < 1350:
        return Const.PR_BG_COLORS.avg;
      case value < 1550:
        return Const.PR_BG_COLORS.good;
      case value < 1750:
        return Const.PR_BG_COLORS.veryGood;
      case value < 2100:
        return Const.PR_BG_COLORS.great;
      case value < 2450:
        return Const.PR_BG_COLORS.unicum;
      case value >= 2450:
        return Const.PR_BG_COLORS.superUnicum;
      default:
        return "";
    }
  }

  export function prText(value: number): string {
    switch (true) {
      case value <= 0:
        return "";
      case value < 750:
        return Const.PR_TEXT_COLORS.bad;
      case value < 1100:
        return Const.PR_TEXT_COLORS.belowAvg;
      case value < 1350:
        return Const.PR_TEXT_COLORS.avg;
      case value < 1550:
        return Const.PR_TEXT_COLORS.good;
      case value < 1750:
        return Const.PR_TEXT_COLORS.veryGood;
      case value < 2100:
        return Const.PR_TEXT_COLORS.great;
      case value < 2450:
        return Const.PR_TEXT_COLORS.unicum;
      case value >= 2450:
        return Const.PR_TEXT_COLORS.superUnicum;
      default:
        return "";
    }
  }

  export function winRate(value: number): string {
    switch (true) {
      case value <= 0:
        return "";
      case value < 47:
        return Const.PR_TEXT_COLORS.bad;
      case value < 50:
        return Const.PR_TEXT_COLORS.belowAvg;
      case value < 52:
        return Const.PR_TEXT_COLORS.avg;
      case value < 54:
        return Const.PR_TEXT_COLORS.good;
      case value < 56:
        return Const.PR_TEXT_COLORS.veryGood;
      case value < 60:
        return Const.PR_TEXT_COLORS.great;
      case value < 65:
        return Const.PR_TEXT_COLORS.unicum;
      case value >= 65:
        return Const.PR_TEXT_COLORS.superUnicum;
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
        return Const.PR_TEXT_COLORS.bad;
      case ratio < 0.8:
        return Const.PR_TEXT_COLORS.belowAvg;
      case ratio < 1.0:
        return Const.PR_TEXT_COLORS.avg;
      case ratio < 1.2:
        return Const.PR_TEXT_COLORS.good;
      case ratio < 1.4:
        return Const.PR_TEXT_COLORS.veryGood;
      case ratio < 1.5:
        return Const.PR_TEXT_COLORS.great;
      case ratio < 1.6:
        return Const.PR_TEXT_COLORS.unicum;
      case ratio >= 1.6:
        return Const.PR_TEXT_COLORS.superUnicum;
      default:
        return "";
    }
  }
}

export default TextColor;
