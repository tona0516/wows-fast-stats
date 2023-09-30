import { DispName } from "src/lib/DispName";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import { RatingFactor } from "src/lib/rating/RatingFactor";
import type { domain } from "wailsjs/go/models";

export const ratingFactors = (): RatingFactor[] => {
  return [
    new RatingFactor("bad", 11, "cv", 0, 750, 0, 0.6, 0, 47),
    new RatingFactor("below_avg", 10, "bb", 750, 1100, 0.6, 0.8, 47, 50),
    new RatingFactor("avg", 9, "bb", 1100, 1350, 0.8, 1.0, 50, 52),
    new RatingFactor("good", 8, "cl", 1350, 1550, 1.0, 1.2, 52, 54),
    new RatingFactor("very_good", 7, "cl", 1550, 1750, 1.2, 1.4, 54, 56),
    new RatingFactor("great", 6, "dd", 1750, 2100, 1.4, 1.5, 56, 60),
    new RatingFactor("unicum", 5, "dd", 2100, 2450, 1.5, 1.6, 60, 65),
    new RatingFactor("super_unicum", 4, "ss", 2450, 9999, 1.6, 10, 65, 100),
  ];
};

export const sampleTeam = (): domain.Team => {
  const avgDamage = 10000;
  const players: domain.Player[] = ratingFactors().map((value, i, _) => {
    const playerInfo: domain.PlayerInfo = {
      id: 1,
      name: "player_name" + i + 1,
      clan: { tag: "TEST" } as domain.Clan,
      is_hidden: false,
      convertValues: function (a: any, classs: any, asMap?: boolean) {
        throw new Error("Function not implemented.");
      },
    };
    const shipInfo: domain.ShipInfo = {
      id: 1,
      name: "Test Ship",
      nation: "japan",
      tier: value.tier,
      type: value.shipType,
      avg_damage: avgDamage,
      is_premium: false,
    };
    const shipStats: domain.ShipStats = {
      battles: 10,
      damage: value.minDamage * avgDamage,
      max_damage: {
        ship_id: 0,
        ship_name: "",
        ship_tier: 0,
        damage: value.minDamage * avgDamage * 1.5,
      },
      win_rate: value.minWin,
      win_survived_rate: 50,
      lose_survived_rate: 50,
      kd_rate: 1,
      kill: 1,
      exp: 1000,
      pr: value.minPR,
      main_battery_hit_rate: 50,
      torpedoes_hit_rate: 5,
      planes_killed: 5,
      convertValues: function (a: any, classs: any, asMap?: boolean) {
        throw new Error("Function not implemented.");
      },
    };
    const overallStats: domain.OverallStats = {
      battles: 10,
      damage: value.minDamage * avgDamage,
      max_damage: {
        ship_id: 0,
        ship_name: "Test Ship",
        ship_tier: 5,
        damage: value.minDamage * avgDamage * 1.5,
      },
      win_rate: value.minWin,
      win_survived_rate: 50,
      lose_survived_rate: 50,
      kd_rate: 1,
      kill: 1,
      exp: 1000,
      pr: value.minPR,
      avg_tier: 5,
      using_ship_type_rate: {
        ss: 20,
        dd: 20,
        cl: 20,
        bb: 20,
        cv: 20,
      } as domain.ShipTypeGroup,
      using_tier_rate: {
        low: 33.3,
        middle: 33.3,
        high: 33.4,
      } as domain.TierGroup,
      convertValues: function (a: any, classs: any, asMap?: boolean) {
        throw new Error("Function not implemented.");
      },
    };

    return {
      player_info: playerInfo,
      ship_info: shipInfo,
      pvp_solo: {
        ship: shipStats,
        overall: overallStats,
      } as domain.PlayerStats,
      pvp_all: {
        ship: shipStats,
        overall: overallStats,
      } as domain.PlayerStats,
      convertValues(a, classs, asMap) {
        throw new Error("Function not implemented.");
      },
    };
  });

  return { players: players } as domain.Team;
};

export const colorDescriptions = (userConfig: domain.UserConfig) => {
  return ratingFactors().map((rating) => {
    return {
      level: { text: DispName.SKILL_LEVELS.get(rating.level) },
      playerName: {
        text: "player_name",
        bgColor: RatingConverterFactory.fromPR(
          rating.minPR,
          userConfig,
        ).bgColorCode(),
      },
      pr: {
        text: `${rating.minPR} ~ ${rating.maxPR}`,
        textColor: RatingConverterFactory.fromDamage(
          rating.minDamage,
          1.0,
          userConfig,
        ).textColorCode(),
      },
      damage: {
        text: `${rating.minDamage}倍 ~ ${rating.maxDamage}倍`,
        textColor: RatingConverterFactory.fromDamage(
          rating.minDamage,
          1.0,
          userConfig,
        ).textColorCode(),
      },
      win: {
        text: `${rating.minWin}% ~ ${rating.maxWin}%`,
        textColor: RatingConverterFactory.fromWinRate(
          rating.minWin,
          userConfig,
        ).textColorCode(),
      },
    };
  });
};
