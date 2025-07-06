<script lang="ts">
  import { formatWithSuffix, toPlayerStats } from "src/lib/util";
  import { storedConfig } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;

  function getPlayerStats(team: data.Team): data.PlayerStats[] {
    return team.players
      .filter((p) => p.player_info.id !== 0 && !p.player_info.is_hidden)
      .map((p) => toPlayerStats(p, $storedConfig.stats_pattern));
  }

  function average(values: number[]): number {
    if (values.length === 0) {
      return 0;
    }
    return values.reduce((a, b) => a + b) / values.length;
  }

  function standardDeviation(values: number[]): number {
    if (values.length === 0) {
      return 0;
    }

    const mean = average(values);

    const squaredDiffs = values.map((value) => {
      const diff = value - mean;
      return diff * diff;
    });

    const variance = squaredDiffs.reduce((a, b) => a + b, 0) / values.length;

    return Math.sqrt(variance);
  }

  function getPRs(team: data.Team): number[] {
    return getPlayerStats(team).map((p) => p.overall.pr);
  }

  function getThreatLevels(team: data.Team): number[] {
    return getPlayerStats(team).map((p) => p.overall.threat_level.modified);
  }

  function getDamages(team: data.Team): number[] {
    return getPlayerStats(team).map((p) => p.overall.damage);
  }

  function getWinRates(team: data.Team): number[] {
    return getPlayerStats(team).map((p) => p.overall.win_rate);
  }

  function getBattles(team: data.Team): number[] {
    return getPlayerStats(team).map((p) => p.overall.battles);
  }

  function maxTeamValue(
    teams: data.Team[],
    func: (team: data.Team) => number[],
  ): number {
    return Math.max(...teams.map((team) => average(func(team))));
  }

  const chartItems = [
    {
      label: "PR",
      func: getPRs,
    },
    {
      label: "戦力評価",
      func: getThreatLevels,
    },
    {
      label: "ダメージ",
      func: getDamages,
    },
    {
      label: "勝率",
      func: getWinRates,
    },
    {
      label: "戦闘数",
      func: getBattles,
    },
  ];

  // Note: https://iro-color.com/colorchart/tone/bright-tone.html
  const chartColors = ["#00A95F", "#EA5532", "#187FC4"];
</script>

<div class="w-3xl">
  <table
    class="charts-css column multiple show-labels data-spacing-10 datasets-spacing-4"
  >
    <tbody class="h-32">
      {#each chartItems as item}
        {@const max = maxTeamValue(battle.teams, item.func)}
        <tr>
          <th scope="row">{item.label}</th>
          {#each battle.teams as team, i}
            {@const avg = average(item.func(team))}
            {@const sd = standardDeviation(item.func(team))}
            <td style="--size: calc({avg}/{max}); --color: {chartColors[i]}">
              <div class="text-center">
                <span class="text-nowrap text-neutral-100 font-semibold"
                  >{formatWithSuffix(avg)}</span
                >
                <span class="text-nowrap text-neutral-100"
                  >(±{formatWithSuffix(sd)})</span
                >
              </div>
            </td>
          {/each}
        </tr>
      {/each}
    </tbody>
  </table>
</div>
