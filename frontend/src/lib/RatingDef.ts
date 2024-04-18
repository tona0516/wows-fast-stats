import clone from "clone";
import type { RatingLevel, ShipType } from "src/lib/types";
import { data } from "wailsjs/go/models";

interface RatingDef {
  level: RatingLevel;
  pr: { min: number; max: number };
  damage: { min: number; max: number };
  winRate: { min: number; max: number };
}

export const RATING_DEFS: RatingDef[] = [
  {
    level: "bad",
    pr: { min: 0, max: 750 },
    damage: { min: 0, max: 0.6 },
    winRate: { min: 0, max: 47 },
  },
  {
    level: "below_avg",
    pr: { min: 750, max: 1100 },
    damage: { min: 0.6, max: 0.8 },
    winRate: { min: 47, max: 50 },
  },
  {
    level: "avg",
    pr: { min: 1100, max: 1350 },
    damage: { min: 0.8, max: 1.0 },
    winRate: { min: 50, max: 52 },
  },
  {
    level: "good",
    pr: { min: 1350, max: 1550 },
    damage: { min: 1.0, max: 1.2 },
    winRate: { min: 52, max: 54 },
  },
  {
    level: "very_good",
    pr: { min: 1550, max: 1750 },
    damage: { min: 1.2, max: 1.4 },
    winRate: { min: 54, max: 56 },
  },
  {
    level: "great",
    pr: { min: 1750, max: 2100 },
    damage: { min: 1.4, max: 1.5 },
    winRate: { min: 56, max: 60 },
  },
  {
    level: "unicum",
    pr: { min: 2100, max: 2450 },
    damage: { min: 1.5, max: 1.6 },
    winRate: { min: 60, max: 65 },
  },
  {
    level: "super_unicum",
    pr: { min: 2450, max: 9999 },
    damage: { min: 1.6, max: 10 },
    winRate: { min: 65, max: 100 },
  },
];

const sampleTeam = (): data.Team => {
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

  const players = RATING_DEFS.map((value, i, _) => {
    const ss = new data.ShipStats();
    ss.battles = 10;
    ss.damage = value.damage.min * SAMPLE_AVG_DAMAGE;
    ss.max_damage = {
      ship_id: 0,
      ship_name: "",
      ship_tier: 0,
      value: value.damage.min * SAMPLE_AVG_DAMAGE * 1.5,
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
    ss.platoon_rate = 3.0;

    const os = new data.OverallStats();
    os.battles = 10;
    os.damage = value.damage.min * SAMPLE_AVG_DAMAGE;
    os.max_damage = {
      ship_id: 0,
      ship_name: "Test Ship",
      ship_tier: 5,
      value: value.damage.min * SAMPLE_AVG_DAMAGE * 1.5,
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
    os.platoon_rate = 3.0;

    const tm = new data.ThreatLevel();
    tm.raw = 19000; // TODO
    tm.modified = 19000; // TODO
    os.threat_level = tm;

    const pi = new data.PlayerInfo();
    pi.id = 1;
    pi.name = "player_name" + i + 1;
    pi.clan = { tag: "TEST", id: 1, hex_color: "" };
    pi.is_hidden = false;

    const si = new data.ShipInfo();
    si.id = 1;
    si.name = "Test Ship";
    si.nation = "japan";
    si.tier = SHIP_INFO_SAMPLES[i].tier;
    si.type = SHIP_INFO_SAMPLES[i].type;
    si.avg_damage = SAMPLE_AVG_DAMAGE;
    si.is_premium = false;

    const pvpSolo = new data.PlayerStats();
    pvpSolo.ship = clone(ss);
    pvpSolo.overall = clone(os);

    const pvpAll = new data.PlayerStats();
    pvpAll.ship = clone(ss);
    pvpAll.overall = clone(os);

    const player = new data.Player();
    player.player_info = pi;
    player.ship_info = si;
    player.pvp_solo = pvpSolo;
    player.pvp_all = pvpAll;

    return player;
  });

  const team = new data.Team();
  team.players = players;

  return team;
};

export const SAMPLE_TEAM = sampleTeam();
