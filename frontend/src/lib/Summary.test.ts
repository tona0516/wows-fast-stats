import { DispName } from "src/lib/DispName";
import { Summary } from "src/lib/Summary";
import type { StatsExtra } from "src/lib/types";
import { data } from "wailsjs/go/models";

const makePlayer = (): data.Player => {
  const pvpSolo = new data.PlayerStats();
  pvpSolo.ship = new data.ShipStats();
  pvpSolo.overall = new data.OverallStats();

  const pvpAll = new data.PlayerStats();
  pvpAll.ship = new data.ShipStats();
  pvpAll.overall = new data.OverallStats();

  const player = new data.Player();
  player.player_info = new data.PlayerInfo();
  player.ship_info = new data.ShipInfo();
  player.pvp_solo = pvpSolo;
  player.pvp_all = pvpAll;

  return player;
};

const makeConfig = (): data.UserConfigV2 => {
  const digit = new data.UCDigit();
  const teamSummary = new data.UCTeamSummary();

  const display = new data.UCDisplay();
  display.ship = new data.UCDisplayShip();
  display.overall = new data.UCDisplayOverall();

  const config = new data.UserConfigV2();
  config.digit = digit;
  config.team_summary = teamSummary;
  config.display = display;

  return config;
};

test("calculate - invalid battle", () => {
  const battles = [
    undefined, // battleが存在しない
    (() => {
      const battle = new data.Battle();
      battle.teams = [new data.Team()];
      return battle;
    })(), // teamが1つ
  ];

  battles.forEach((it) => {
    const actual = Summary.calculate(it, new Set(), new data.UserConfigV2());
    expect(actual).toBeUndefined();
  });
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

  const friendTeam = new data.Team();
  friendTeam.players = [friend1, friend2];
  const enemyTeam = new data.Team();
  enemyTeam.players = [enemy1, enemy2];

  const battle = new data.Battle();
  battle.teams = [friendTeam, enemyTeam];

  const config = makeConfig();
  config.digit.pr = 0;
  config.digit.damage = 1;
  config.digit.win_rate = 2;
  config.team_summary.min_ship_battles = 20;
  config.stats_pattern = extra;

  const summary = Summary.calculate(
    battle,
    new Set([enemy2.player_info.id]),
    config,
  );

  expect(summary?.values.get("all")?.friends).toEqual([
    friend1[extra].ship.pr.toFixed(config.digit.pr),
    friend1[extra].ship.damage.toFixed(config.digit.damage),
    friend1[extra].ship.win_rate.toFixed(config.digit.win_rate),
    "-",
    "-",
    "-",
  ]);
  expect(summary?.values.get("all")?.enemies).toEqual([
    enemy1[extra].ship.pr.toFixed(config.digit.pr),
    enemy1[extra].ship.damage.toFixed(config.digit.damage),
    enemy1[extra].ship.win_rate.toFixed(config.digit.win_rate),
    "-",
    "-",
    "-",
  ]);
  expect(summary?.values.get("all")?.diffs).toEqual([
    {
      colorCode: "",
      diff: Math.abs(friend1[extra].ship.pr - enemy1[extra].ship.pr).toFixed(
        config.digit.pr,
      ),
    },
    {
      colorCode: "#fc4e32",
      diff: `-${Math.abs(
        friend1[extra].ship.damage - enemy1[extra].ship.damage,
      ).toFixed(config.digit.damage)}`,
    },
    {
      colorCode: "#99d02b",
      diff: `+${Math.abs(
        friend1[extra].ship.win_rate - enemy1[extra].ship.win_rate,
      ).toFixed(config.digit.win_rate)}`,
    },
    { colorCode: "", diff: "-" },
    { colorCode: "", diff: "-" },
    { colorCode: "", diff: "-" },
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

  const friendTeam = new data.Team();
  friendTeam.players = friends;
  const enemyTeam = new data.Team();
  enemyTeam.players = enemies;

  const battle = new data.Battle();
  battle.teams = [friendTeam, enemyTeam];

  const config = makeConfig();
  config.digit.pr = 0;
  config.digit.damage = 1;
  config.digit.win_rate = 2;
  config.team_summary.min_overall_battles = 1;
  config.stats_pattern = extra;

  const summary = Summary.calculate(battle, new Set(), config);

  shipTypes.forEach((shipType) => {
    expect(summary?.values.get(shipType)?.friends).toEqual([
      "-",
      "-",
      "-",
      pr.toFixed(config.digit.pr),
      damage.toFixed(config.digit.damage),
      winRate.toFixed(config.digit.win_rate),
    ]);
  });

  shipTypes.forEach((shipType) => {
    expect(summary?.values.get(shipType)?.diffs).toEqual([
      { colorCode: "", diff: "-" },
      { colorCode: "", diff: "-" },
      { colorCode: "", diff: "-" },
      { colorCode: "", diff: "-" },
      { colorCode: "", diff: "-" },
      { colorCode: "", diff: "-" },
    ]);
  });
});
