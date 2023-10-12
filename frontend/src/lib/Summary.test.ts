import { Summary } from "src/lib/Summary";
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

const makeTeam = (players: domain.Player[]): domain.Team => {
  const team = new domain.Team();
  team.players = players;

  return team;
};

test("undefined", () => {
  expect(
    Summary.calculate(undefined, [], new domain.UserConfig()),
  ).toBeUndefined();
});

test("normal", () => {
  const friend1 = makePlayer();
  friend1.player_info.id = 1;
  friend1.pvp_all.ship.battles = 100;
  friend1.pvp_all.ship.pr = 1000;
  friend1.pvp_all.ship.damage = 10000;
  friend1.pvp_all.ship.win_rate = 50;

  const friend2 = makePlayer();
  friend2.player_info.id = 2;
  // Note: under min_ship_battles
  friend2.pvp_all.ship.battles = 19;

  const enemy1 = makePlayer();
  enemy1.player_info.id = 11;
  enemy1.pvp_all.ship.battles = 100;
  enemy1.pvp_all.ship.pr = friend1.pvp_all.ship.pr;
  enemy1.pvp_all.ship.damage = friend1.pvp_all.ship.damage + 10;
  enemy1.pvp_all.ship.win_rate = friend1.pvp_all.ship.win_rate - 10;

  const enemy2 = makePlayer();
  // Note: into excluded players
  enemy2.player_info.id = 12;
  enemy2.pvp_all.ship.battles = 100;

  const battle = new domain.Battle();
  battle.teams = [makeTeam([friend1, friend2]), makeTeam([enemy1, enemy2])];

  const customDigit = new domain.CustomDigit();
  customDigit.pr = 0;
  customDigit.damage = 1;
  customDigit.win_rate = 2;

  const teamAvg = new domain.TeamAverage();
  teamAvg.min_ship_battles = 20;

  const config = new domain.UserConfig();
  config.custom_digit = customDigit;
  config.team_average = teamAvg;
  config.stats_pattern = "pvp_all";

  const summary = Summary.calculate(battle, [enemy2.player_info.id], config);

  expect(summary?.values.friends).toEqual([
    friend1.pvp_all.ship.pr.toFixed(customDigit.pr),
    friend1.pvp_all.ship.damage.toFixed(customDigit.damage),
    friend1.pvp_all.ship.win_rate.toFixed(customDigit.win_rate),
    (0).toFixed(customDigit.pr),
    (0).toFixed(customDigit.damage),
    (0).toFixed(customDigit.win_rate),
  ]);
  expect(summary?.values.enemies).toEqual([
    enemy1.pvp_all.ship.pr.toFixed(customDigit.pr),
    enemy1.pvp_all.ship.damage.toFixed(customDigit.damage),
    enemy1.pvp_all.ship.win_rate.toFixed(customDigit.win_rate),
    (0).toFixed(customDigit.pr),
    (0).toFixed(customDigit.damage),
    (0).toFixed(customDigit.win_rate),
  ]);
  expect(summary?.values.diffs).toEqual([
    {
      colorCode: "",
      diff: Math.abs(friend1.pvp_all.ship.pr - enemy1.pvp_all.ship.pr).toFixed(
        customDigit.pr,
      ),
    },
    {
      colorCode: "#fc4e32",
      diff: `-${Math.abs(
        friend1.pvp_all.ship.damage - enemy1.pvp_all.ship.damage,
      ).toFixed(customDigit.damage)}`,
    },
    {
      colorCode: "#99d02b",
      diff: `+${Math.abs(
        friend1.pvp_all.ship.win_rate - enemy1.pvp_all.ship.win_rate,
      ).toFixed(customDigit.win_rate)}`,
    },
    { colorCode: "", diff: (0).toFixed(customDigit.pr) },
    { colorCode: "", diff: (0).toFixed(customDigit.damage) },
    { colorCode: "", diff: (0).toFixed(customDigit.win_rate) },
  ]);
});
