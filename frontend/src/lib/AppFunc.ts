import chroma from "chroma-js";
import { geometricMean } from "simple-statistics";
import { data } from "wailsjs/go/models";
import { AppConst } from "./AppConst";
import { LocalStorage } from "./LocalStorage";
import type {
  ColorPair,
  RatingInfo,
  RowPattern,
  ShipType,
  StatsExtra,
  TeamThreatLevel,
  ThreatLevelInfo,
} from "./types";

export namespace AppFunc {
  export const isLighter = (): boolean => {
    const current = LocalStorage.instance.getTheme();
    return AppConst.LIGHTER_THEMES.includes(current);
  };

  export const clanNumbersURL = (clanID: number): string =>
    `${AppConst.NUMBERS_URL}clan/${clanID},/`;

  export const playerNumbersURL = (
    accountID: number,
    accountName: string,
  ): string => `${AppConst.NUMBERS_URL}player/${accountID},${accountName}/`;

  export const shipNumbersURL = (shipID: number): string =>
    `${AppConst.NUMBERS_URL}ship/${shipID},/`;

  export const getThreatLevel = (
    score: number,
  ): ThreatLevelInfo | undefined => {
    return AppConst.THREAT_LEVELS_INFO.findLast((it) => score >= it.threshold);
  };

  export const getRatingfromPR = (value: number): RatingInfo | undefined => {
    return AppConst.RATING_INFO.findLast((it) => value >= it.prThreshold);
  };

  export const getRatingfromWinRate = (
    value: number,
  ): RatingInfo | undefined => {
    return AppConst.RATING_INFO.findLast((it) => value >= it.winRateThreshold);
  };

  export const getRatingfromDamage = (
    value: number,
    expected: number,
  ): RatingInfo | undefined => {
    if (expected === 0) {
      return undefined;
    }
    const ratio = value / expected;
    return AppConst.RATING_INFO.findLast((it) => ratio >= it.damageThreshold);
  };

  export const getTeamThreatLevels = (
    battle: data.Battle | undefined,
    statsExtra: StatsExtra,
  ): TeamThreatLevel[] => {
    if (!battle || !battle.teams) {
      return [];
    }

    return battle.teams.map((team) => {
      const players = team.players;
      const values = players
        .filter((player) => {
          const id = player.player_info.id;
          return !(id === 0 || player.player_info.is_hidden);
        })
        .map((player) => player[statsExtra].overall.threat_level.modified);
      const maxScore = Math.max(...values);

      const average = geometricMean(values);

      return {
        average: average,
        dissociationDegree: (maxScore / average - 1) * 100,
        accuracy: Math.round((values.length / players.length) * 100),
      };
    });
  };

  export const getFixedColorPair = (colorCode: string): ColorPair => {
    const chromaColor = chroma(colorCode);
    const colorRate = 1.5;
    const brighten = chromaColor.brighten(colorRate).hex();
    const darken = chromaColor.darken(colorRate).hex();

    if (AppFunc.isLighter()) {
      return { text: darken, background: brighten };
    }
    return { text: brighten, background: darken };
  };

  export const toTierString = (value: number): string => {
    if (value === 11) return "★";
    return AppConst.ROMAN_NUMERALS[value] ?? "";
  };

  export const toShipType = (type: string): type is ShipType => {
    return Object.keys(new data.ShipTypeGroup()).includes(type);
  };

  const formatNumber = (value: number, digit: number): string => {
    return value.toLocaleString("ja-JP", {
      minimumFractionDigits: digit,
      maximumFractionDigits: digit,
    });
  };

  export const formatWithSuffix = (num: number): string => {
    const G = 1_000_000_000;
    const M = 1_000_000;
    const K = 1_000;

    if (num >= G) {
      return `${formatNumber(num / G, 1)}G`;
    }

    if (num >= M) {
      return `${formatNumber(num / M, 1)}M`;
    }

    if (num >= K) {
      return `${formatNumber(num / K, 1)}K`;
    }

    return formatNumber(num, 0);
  };

  export const getRowPattern = (
    player: data.Player,
    statsExtra: StatsExtra,
    shipColumnCount: number,
    overallColumnCount: number,
  ): RowPattern => {
    if (shipColumnCount + overallColumnCount === 0) {
      return "no_column";
    }

    if (player.player_info.is_hidden === true) {
      return "private";
    }

    const stats = player[statsExtra];
    if (player.player_info.id === 0 || stats.overall.battles === 0) {
      return "no_stats";
    }

    if (stats.ship.battles === 0 && shipColumnCount > 0) {
      return "no_ship_stats";
    }

    return "full";
  };

  export const getColumnText = (pattern: RowPattern): string => {
    switch (pattern) {
      case "private":
        return "PRIVAYE";
      case "no_stats":
        return "N/A";
      case "no_ship_stats":
        return "N/A";
      default:
        return "";
    }
  };
}
