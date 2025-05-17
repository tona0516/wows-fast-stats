import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
import { data } from "wailsjs/go/models";

test("fromBattle - 異常系", () => {
  expect(
    TeamThreatLevel.fromBattle(undefined, new Set(), "pvp_all"),
  ).toBeUndefined();
});

test("fromBattle - 正常系", () => {
  const values: { threatLevel: number; isHidden: boolean; id: number }[] = [
    { threatLevel: 10000, isHidden: false, id: 1 },
    { threatLevel: 20000, isHidden: false, id: 2 },
    { threatLevel: 10000, isHidden: true, id: 3 }, // private
    { threatLevel: 10000, isHidden: false, id: 0 }, // npc
  ];

  const battle: data.Battle = {
    teams: [
      {
        players: values.map((value) => {
          return {
            player_info: { is_hidden: value.isHidden },
            pvp_all: {
              overall: {
                threat_level: {
                  modified: value.threatLevel,
                },
              },
            },
          } as data.Player;
        }),
        convertValues: () => {},
      },
    ],
    convertValues: () => {},
    meta: new data.Meta(),
  };

  // biome-ignore lint/style/noNonNullAssertion: <explanation>
  const actual = TeamThreatLevel.fromBattle(battle, new Set(), "pvp_all")![0];
  const expected = new TeamThreatLevel(12599, 59, 75);

  expect(actual.average.toFixed()).toBe(expected.average.toFixed());
  expect(actual.dissociationDegree.toFixed()).toBe(
    expected.dissociationDegree.toFixed(),
  );
  expect(actual.accuracy.toFixed()).toBe(expected.accuracy.toFixed());
});
