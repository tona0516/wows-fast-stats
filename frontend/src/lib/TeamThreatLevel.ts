import type { OptionalTeamThreatLevels, StatsExtra } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class TeamThreatLevel {
	constructor(
		readonly average: number,
		readonly dissociationDegree: number,
		readonly accuracy: number,
	) {}

	static fromBattle = (
		battle: data.Battle,
		statsExtra: StatsExtra,
	): Map<string, TeamThreatLevel> => {
		const result = new Map<string, TeamThreatLevel>();

		for (const [i, team] of battle.teams.entries()) {
			const players = team.players;
			const values = players
				.filter((player) => {
					const id = player.player_info.id;
					return !(id === 0 || player.player_info.is_hidden);
				})
				.map((player) => player[statsExtra].overall.threat_level.modified);
			const maxScore = Math.max(...values);
			const average = calcGeometricMean(values);

			const teamName = i === 0 ? "味方" : "敵";

			result.set(
				teamName,
				new TeamThreatLevel(
					average,
					(maxScore / average - 1) * 100,
					Math.round((values.length / players.length) * 100),
				),
			);
		}

		return result;
	};
}

const calcGeometricMean = (values: number[]): number => {
	if (values.length === 0) {
		return 0;
	}

	const productScore = values.reduce((a, b) => a * b, 1);

	return productScore ** (1 / values.length);
};
