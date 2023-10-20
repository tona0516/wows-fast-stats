import { DispName } from "src/lib/DispName";
import { Summary } from "src/lib/Summary";
import type { StatsExtra } from "src/lib/types";
import { domain } from "wailsjs/go/models";

const makePlayer = (): domain.Player => {
  const pvpSolo = new domain.PlayerStats();
  pvpSolo.ship = new domain.ShipStats();
  pvpSolo.overall = new domain.OverallStats();

  const pvpAll = new domain.PlayerStats();
  pvpAll.ship = new domain.ShipStats();
  pvpAll.overall = new domain.OverallStats();

  const player = new domain.Player();
  player.player_info = new domain.PlayerInfo();
  player.ship_info = new domain.ShipInfo();
  player.pvp_solo = pvpSolo;
  player.pvp_all = pvpAll;

  return player;
};

test("undefined", () => {
  expect(
    Summary.calculate(undefined, [], new domain.UserConfig()),
  ).toBeUndefined();
});

test("calculate - all types, ship, pvp_all, excluded player", () => {
  const extra: StatsExtra = "pvp_all";

  const friend1 = makePlayer();
  friend1.player_info.id = 1;
  friend1[extra].ship.battles = 100;
  friend1[extra].ship.pr = 1000;
  friend1[extra].ship.damage = 10000;
  friend1[extra].ship.win_rate = 50;

  const friend2 = makePlayer();
  friend2.player_info.id = 2;
  // Note: under min_ship_battles
  friend2.pvp_all.ship.battles = 19;

  const enemy1 = makePlayer();
  enemy1.player_info.id = 11;
  enemy1[extra].ship.battles = 100;
  enemy1[extra].ship.pr = friend1[extra].ship.pr;
  enemy1[extra].ship.damage = friend1[extra].ship.damage + 10;
  enemy1[extra].ship.win_rate = friend1[extra].ship.win_rate - 10;

  const enemy2 = makePlayer();
  // Note: into excluded players
  enemy2.player_info.id = 12;
  enemy2[extra].ship.battles = 100;

  const friendTeam = new domain.Team();
  friendTeam.players = [friend1, friend2];
  const enemyTeam = new domain.Team();
  enemyTeam.players = [enemy1, enemy2];

  const battle = new domain.Battle();
  battle.teams = [friendTeam, enemyTeam];

  const customDigit = new domain.CustomDigit();
  customDigit.pr = 0;
  customDigit.damage = 1;
  customDigit.win_rate = 2;

  const teamAvg = new domain.TeamAverage();
  teamAvg.min_ship_battles = 20;

  const config = new domain.UserConfig();
  config.custom_digit = customDigit;
  config.team_average = teamAvg;
  config.stats_pattern = extra;

  const summary = Summary.calculate(battle, [enemy2.player_info.id], config);

  expect(summary?.values.get("all")?.friends).toEqual([
    friend1[extra].ship.pr.toFixed(customDigit.pr),
    friend1[extra].ship.damage.toFixed(customDigit.damage),
    friend1[extra].ship.win_rate.toFixed(customDigit.win_rate),
    (0).toFixed(customDigit.pr),
    (0).toFixed(customDigit.damage),
    (0).toFixed(customDigit.win_rate),
  ]);
  expect(summary?.values.get("all")?.enemies).toEqual([
    enemy1[extra].ship.pr.toFixed(customDigit.pr),
    enemy1[extra].ship.damage.toFixed(customDigit.damage),
    enemy1[extra].ship.win_rate.toFixed(customDigit.win_rate),
    (0).toFixed(customDigit.pr),
    (0).toFixed(customDigit.damage),
    (0).toFixed(customDigit.win_rate),
  ]);
  expect(summary?.values.get("all")?.diffs).toEqual([
    {
      colorCode: "",
      diff: Math.abs(friend1[extra].ship.pr - enemy1[extra].ship.pr).toFixed(
        customDigit.pr,
      ),
    },
    {
      colorCode: "#fc4e32",
      diff: `-${Math.abs(
        friend1[extra].ship.damage - enemy1[extra].ship.damage,
      ).toFixed(customDigit.damage)}`,
    },
    {
      colorCode: "#99d02b",
      diff: `+${Math.abs(
        friend1[extra].ship.win_rate - enemy1[extra].ship.win_rate,
      ).toFixed(customDigit.win_rate)}`,
    },
    { colorCode: "", diff: (0).toFixed(customDigit.pr) },
    { colorCode: "", diff: (0).toFixed(customDigit.damage) },
    { colorCode: "", diff: (0).toFixed(customDigit.win_rate) },
  ]);
});

test("calculate - each ship type, overall, pvp_solo", () => {
  const extra: StatsExtra = "pvp_solo";
  const shipTypes = [...DispName.SHIP_TYPES.keys()];
  const battles = 100,
    pr = 1000,
    damage = 10000,
    winRate = 50;

  const friends = shipTypes.map((it) => {
    const friend = makePlayer();
    friend.ship_info.type = it;
    friend.player_info.id = 1;
    friend[extra].overall.battles = battles;
    friend[extra].overall.pr = pr;
    friend[extra].overall.damage = damage;
    friend[extra].overall.win_rate = winRate;

    return friend;
  });

  const enemies = shipTypes.map((it) => {
    const enemy = makePlayer();
    enemy.ship_info.type = it;
    enemy.player_info.id = 1;
    enemy[extra].overall.battles = 0;
    enemy[extra].overall.pr = 0;
    enemy[extra].overall.damage = 0;
    enemy[extra].overall.win_rate = 0;

    return enemy;
  });

  const friendTeam = new domain.Team();
  friendTeam.players = friends;
  const enemyTeam = new domain.Team();
  enemyTeam.players = enemies;

  const battle = new domain.Battle();
  battle.teams = [friendTeam, enemyTeam];

  const customDigit = new domain.CustomDigit();
  customDigit.pr = 0;
  customDigit.damage = 1;
  customDigit.win_rate = 2;

  const teamAvg = new domain.TeamAverage();
  teamAvg.min_overall_battles = 1;

  const config = new domain.UserConfig();
  config.custom_digit = customDigit;
  config.team_average = teamAvg;
  config.stats_pattern = extra;

  const summary = Summary.calculate(battle, [], config);

  shipTypes.forEach((shipType) => {
    expect(summary?.values.get(shipType)?.friends).toEqual([
      (0).toFixed(customDigit.pr),
      (0).toFixed(customDigit.damage),
      (0).toFixed(customDigit.win_rate),
      pr.toFixed(customDigit.pr),
      damage.toFixed(customDigit.damage),
      winRate.toFixed(customDigit.win_rate),
    ]);
  });

  shipTypes.forEach((shipType) => {
    expect(summary?.values.get(shipType)?.diffs).toEqual([
      { colorCode: "", diff: (0).toFixed(customDigit.pr) },
      { colorCode: "", diff: (0).toFixed(customDigit.damage) },
      { colorCode: "", diff: (0).toFixed(customDigit.win_rate) },
      {
        colorCode: "#99d02b",
        diff: `+${pr.toFixed(customDigit.pr)}`,
      },
      {
        colorCode: "#99d02b",
        diff: `+${damage.toFixed(customDigit.damage)}`,
      },
      {
        colorCode: "#99d02b",
        diff: `+${winRate.toFixed(customDigit.win_rate)}`,
      },
    ]);
  });
});
