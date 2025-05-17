import type {
  OptionalBattle,
  OptionalTeamThreatLevels,
  StatsExtra,
} from "src/lib/types";

export class TeamThreatLevel {
  constructor(
    readonly average: number,
    readonly dissociationDegree: number,
    readonly accuracy: number,
  ) {}

  static fromBattle = (
    battle: OptionalBattle,
    excludedPlayers: Set<number>,
    statsExtra: StatsExtra,
  ): OptionalTeamThreatLevels => {
    if (!battle) {
      return undefined;
    }

    return battle.teams.map((team) => {
      const players = team.players;
      const values = players
        .filter((player) => {
          const id = player.player_info.id;
          return !(
            id === 0 ||
            excludedPlayers.has(id) ||
            player.player_info.is_hidden
          );
        })
        .map((player) => player[statsExtra].overall.threat_level.modified);
      const maxScore = Math.max(...values);
      const average = calcGeometricMean(values);

      return new TeamThreatLevel(
        average,
        (maxScore / average - 1) * 100,
        Math.round((values.length / players.length) * 100),
      );
    });
  };
}

const calcGeometricMean = (values: number[]): number => {
  if (values.length === 0) {
    return 0;
  }

  const productScore = values.reduce((a, b) => a * b, 1);

  return productScore ** (1 / values.length);
};
