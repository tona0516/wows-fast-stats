import type { Rating, ShipType } from "src/lib/types";
import type { domain } from "wailsjs/go/models";

interface RatingRange {
  min: number;
  max: number;
}

class RatingFactor {
  constructor(
    public level: Rating,
    public pr: RatingRange,
    public damage: RatingRange,
    public winRate: RatingRange,
  ) {}
}

export const RATING_FACTORS: RatingFactor[] = [
  new RatingFactor(
    "bad",
    { min: 0, max: 750 },
    { min: 0, max: 0.6 },
    { min: 0, max: 47 },
  ),
  new RatingFactor(
    "below_avg",
    { min: 750, max: 1100 },
    { min: 0.6, max: 0.8 },
    { min: 47, max: 50 },
  ),
  new RatingFactor(
    "avg",
    { min: 1100, max: 1350 },
    { min: 0.8, max: 1.0 },
    { min: 50, max: 52 },
  ),
  new RatingFactor(
    "good",
    { min: 1350, max: 1550 },
    { min: 1.0, max: 1.2 },
    { min: 52, max: 54 },
  ),
  new RatingFactor(
    "very_good",
    { min: 1550, max: 1750 },
    { min: 1.2, max: 1.4 },
    { min: 54, max: 56 },
  ),
  new RatingFactor(
    "great",
    { min: 1750, max: 2100 },
    { min: 1.4, max: 1.5 },
    { min: 56, max: 60 },
  ),
  new RatingFactor(
    "unicum",
    { min: 2100, max: 2450 },
    { min: 1.5, max: 1.6 },
    { min: 60, max: 65 },
  ),
  new RatingFactor(
    "super_unicum",
    { min: 2450, max: 9999 },
    { min: 1.6, max: 10 },
    { min: 65, max: 100 },
  ),
];

const SHIP_INFO_FOR_SAMPLE: { shipType: ShipType; tier: number }[] = [
  { shipType: "cv", tier: 11 },
  { shipType: "bb", tier: 10 },
  { shipType: "bb", tier: 9 },
  { shipType: "cl", tier: 8 },
  { shipType: "cl", tier: 7 },
  { shipType: "dd", tier: 6 },
  { shipType: "dd", tier: 5 },
  { shipType: "ss", tier: 4 },
];

const convertValues = (_a: any, _classs: any, _asMap?: boolean) => {
  throw new Error("Function not implemented.");
};

export const sampleTeam = (): domain.Team => {
  const SAMPLE_AVG_DAMAGE = 10000;

  const players: domain.Player[] = RATING_FACTORS.map((value, i, _) => {
    const shipStats: domain.ShipStats = {
      battles: 10,
      damage: value.damage.min * SAMPLE_AVG_DAMAGE,
      max_damage: {
        ship_id: 0,
        ship_name: "",
        ship_tier: 0,
        damage: value.damage.min * SAMPLE_AVG_DAMAGE * 1.5,
      },
      win_rate: value.winRate.min,
      win_survived_rate: 50,
      lose_survived_rate: 50,
      kd_rate: 1,
      kill: 1,
      exp: 1000,
      pr: value.pr.min,
      main_battery_hit_rate: 50,
      torpedoes_hit_rate: 5,
      planes_killed: 5,
      convertValues,
    };
    const overallStats: domain.OverallStats = {
      battles: 10,
      damage: value.damage.min * SAMPLE_AVG_DAMAGE,
      max_damage: {
        ship_id: 0,
        ship_name: "Test Ship",
        ship_tier: 5,
        damage: value.damage.min * SAMPLE_AVG_DAMAGE * 1.5,
      },
      win_rate: value.winRate.min,
      win_survived_rate: 50,
      lose_survived_rate: 50,
      kd_rate: 1,
      kill: 1,
      exp: 1000,
      pr: value.pr.min,
      avg_tier: 5,
      using_ship_type_rate: {
        ss: 20,
        dd: 20,
        cl: 20,
        bb: 20,
        cv: 20,
      },
      using_tier_rate: {
        low: 33.3,
        middle: 33.3,
        high: 33.4,
      },
      convertValues,
    };

    return {
      player_info: {
        id: 1,
        name: "player_name" + i + 1,
        clan: { tag: "TEST" } as domain.Clan,
        is_hidden: false,
        convertValues,
      },
      ship_info: {
        id: 1,
        name: "Test Ship",
        nation: "japan",
        tier: SHIP_INFO_FOR_SAMPLE[i].tier,
        type: SHIP_INFO_FOR_SAMPLE[i].shipType,
        avg_damage: SAMPLE_AVG_DAMAGE,
        is_premium: false,
      },
      pvp_solo: {
        ship: shipStats,
        overall: overallStats,
        convertValues,
      },
      pvp_all: {
        ship: shipStats,
        overall: overallStats,
        convertValues,
      },
      convertValues,
    };
  });

  return { name: "", players: players, convertValues };
};
