import clone from "clone";
import { RatingAdapter } from "src/lib/rating/RatingColor";
import type { Rating, ShipType } from "src/lib/types";
import { model } from "wailsjs/go/models";

class RatingFactor {
  constructor(
    public level: Rating,
    public pr: { min: number; max: number },
    public damage: { min: number; max: number },
    public winRate: { min: number; max: number },
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

const getSampleTeam = (): model.Team => {
  const SAMPLE_AVG_DAMAGE = 10000;
  const SHIP_INFO_SAMPLES: { type: ShipType; tier: number }[] = [
    { type: "cv", tier: 11 },
    { type: "bb", tier: 10 },
    { type: "bb", tier: 9 },
    { type: "cl", tier: 8 },
    { type: "cl", tier: 7 },
    { type: "dd", tier: 6 },
    { type: "dd", tier: 5 },
    { type: "ss", tier: 4 },
  ];

  const players = RATING_FACTORS.map((value, i, _) => {
    const ss = new model.ShipStats();
    ss.battles = 10;
    ss.damage = value.damage.min * SAMPLE_AVG_DAMAGE;
    ss.max_damage = {
      ship_id: 0,
      ship_name: "",
      ship_tier: 0,
      damage: value.damage.min * SAMPLE_AVG_DAMAGE * 1.5,
    };
    ss.win_rate = value.winRate.min;
    ss.win_survived_rate = 50;
    ss.lose_survived_rate = 50;
    ss.kd_rate = 1;
    ss.kill = 1;
    ss.exp = 1000;
    ss.pr = value.pr.min;
    ss.main_battery_hit_rate = 50;
    ss.torpedoes_hit_rate = 5;
    ss.planes_killed = 5;
    ss.platoon_rate = 3.00;

    const os = new model.OverallStats();
    os.battles = 10;
    os.damage = value.damage.min * SAMPLE_AVG_DAMAGE;
    os.max_damage = {
      ship_id: 0,
      ship_name: "Test Ship",
      ship_tier: 5,
      damage: value.damage.min * SAMPLE_AVG_DAMAGE * 1.5,
    };
    os.win_rate = value.winRate.min;
    os.win_survived_rate = 50;
    os.lose_survived_rate = 50;
    os.kd_rate = 1;
    os.kill = 1;
    os.exp = 1000;
    os.pr = value.pr.min;
    os.avg_tier = 5;
    os.using_ship_type_rate = {
      ss: 20,
      dd: 20,
      cl: 20,
      bb: 20,
      cv: 20,
    };
    os.using_tier_rate = {
      low: 33.3,
      middle: 33.3,
      high: 33.4,
    };
    os.platoon_rate = 3.00;

    const pi = new model.PlayerInfo();
    pi.id = 1;
    pi.name = "player_name" + i + 1;
    pi.clan = { tag: "TEST", id: 1, hex_color: "" };
    pi.is_hidden = false;

    const si = new model.ShipInfo();
    si.id = 1;
    si.name = "Test Ship";
    si.nation = "japan";
    si.tier = SHIP_INFO_SAMPLES[i].tier;
    si.type = SHIP_INFO_SAMPLES[i].type;
    si.avg_damage = SAMPLE_AVG_DAMAGE;
    si.is_premium = false;

    const pvpSolo = new model.PlayerStats();
    pvpSolo.ship = clone(ss);
    pvpSolo.overall = clone(os);

    const pvpAll = new model.PlayerStats();
    pvpAll.ship = clone(ss);
    pvpAll.overall = clone(os);

    const player = new model.Player();
    player.player_info = pi;
    player.ship_info = si;
    player.pvp_solo = pvpSolo;
    player.pvp_all = pvpAll;

    return player;
  });

  const team = new model.Team();
  team.players = players;

  return team;
};

export const SAMPLE_TEAM = getSampleTeam();

export namespace RatingColorFactory {
  export const fromPR = (
    value: number,
    config: model.UserConfig,
  ): RatingAdapter => {
    const rf = RATING_FACTORS.findLast(
      (it) => value >= 0 && value >= it.pr.min,
    );
    return new RatingAdapter(rf?.level, config);
  };

  export const fromDamage = (
    value: number,
    expected: number,
    config: model.UserConfig,
  ): RatingAdapter => {
    if (expected === 0) {
      return new RatingAdapter(undefined, config);
    }

    const ratio = value / expected;
    const rf = RATING_FACTORS.findLast((it) => ratio >= it.damage.min);
    return new RatingAdapter(rf?.level, config);
  };

  export const fromWinRate = (
    value: number,
    config: model.UserConfig,
  ): RatingAdapter => {
    const rf = RATING_FACTORS.findLast((it) => value >= it.winRate.min);
    return new RatingAdapter(rf?.level, config);
  };
}
