<script lang="ts">
  import { max, mean, standardDeviation } from "simple-statistics";
  import type { GetStatsFunction } from "src/lib/types";
  import { formatWithSuffix, toPlayerStats } from "src/lib/util";
  import { storedConfig } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;

  // Note: https://iro-color.com/colorchart/tone/bright-tone.html
  const CHART_COLORS = ["#00A95F", "#EA5532"];
  const CHART_INFO: { label: string; func: GetStatsFunction }[] = [
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

  function getMaxValueInAllPlayers(
    teams: data.Team[],
    func: GetStatsFunction,
  ): number {
    let tempMax = 1;
    for (const team of teams) {
      const values = getValuesInTeam(team, func);
      const avg = calculateMean(values);
      if (avg < tempMax) {
        continue;
      }

      tempMax = avg;
    }

    return tempMax;
  }

  function getValuesInTeam(team: data.Team, func: GetStatsFunction): number[] {
    return team.players
      .filter((p) => p.player_info.id !== 0 && !p.player_info.is_hidden)
      .map((p) => toPlayerStats(p, $storedConfig.stats_pattern))
      .map(func);
  }

  function calculateMean(values: number[]): number {
    return values.length === 0 ? 0 : mean(values);
  }

  function calculateStandardDeviation(values: number[]): number {
    return values.length === 0 ? 0 : standardDeviation(values);
  }
</script>

<div class="w-3xl">
  <table
    class="charts-css column multiple show-labels data-spacing-4 datasets-spacing-4"
  >
    <tbody class="h-32">
      {#each CHART_INFO as item}
        {@const max = getMaxValueInAllPlayers(battle.teams, item.func)}
        <tr>
          <th scope="row">{item.label}</th>
          {#each battle.teams as team, i}
            {@const values = getValuesInTeam(team, item.func)}
            {@const mn = calculateMean(values)}
            {@const sd = calculateStandardDeviation(values)}

            <td style="--size: calc({mn}/{max}); --color: {CHART_COLORS[i]}">
              <div
                class="text-center {mn / max <= 0 &&
                  'absolute -translate-y-full'}"
              >
                <span class="text-nowrap text-neutral-100 font-semibold">
                  {formatWithSuffix(mn)}
                </span>
                <span class="text-nowrap text-neutral-100">
                  (±{formatWithSuffix(sd)})
                </span>
              </div>
            </td>
          {/each}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .text-above-bar {
    position: absolute;
    transform: translateY(-100%);
  }
</style>
