<script lang="ts">
  import { mean, standardDeviation } from "simple-statistics";
  import { formatWithSuffix, toPlayerStatsValues } from "src/lib/util";
  import { storedConfig } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;

  function maxTeamValue(teams: data.Team[], func: GetStatsFunction): number {
    return Math.max(...teams.map((team) => mean(getPlayers(team).map(func))));
  }

  const getPlayers = (team: data.Team): data.PlayerStats[] => {
    return toPlayerStatsValues(team, $storedConfig.stats_pattern);
  };

  type GetStatsFunction = (ps: data.PlayerStats) => number;

  interface ChartInfo {
    label: string;
    func: GetStatsFunction;
  }
  const chartInfos: ChartInfo[] = [
    {
      label: "PR",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.pr;
      },
    },
    {
      label: "戦力評価",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.threat_level.modified;
      },
    },
    {
      label: "ダメージ",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.damage;
      },
    },
    {
      label: "勝率",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.win_rate;
      },
    },
    {
      label: "戦闘数",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.battles;
      },
    },
  ];

  // Note: https://iro-color.com/colorchart/tone/bright-tone.html
  const chartColors = ["#00A95F", "#EA5532"];
</script>

<div class="w-3xl">
  <table
    class="charts-css column multiple show-labels data-spacing-4 datasets-spacing-4"
  >
    <tbody class="h-32">
      {#each chartInfos as item}
        {@const max = maxTeamValue(battle.teams, item.func)}
        <tr>
          <th scope="row">{item.label}</th>
          {#each battle.teams as team, i}
            {@const values = getPlayers(team)}
            {@const mn = mean(values.map(item.func))}
            {@const sd = standardDeviation(values.map(item.func))}

            <td style="--size: calc({mn}/{max}); --color: {chartColors[i]}">
              <div class="text-center">
                <span class="text-nowrap text-neutral-100 font-semibold"
                  >{formatWithSuffix(mn)}</span
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
