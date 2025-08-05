<script lang="ts">
  import { mean, standardDeviation } from "simple-statistics";
  import type { GetStatsFunction } from "src/lib/types";
  import { formatWithSuffix, toPlayerStats } from "src/lib/util";
  import { storedConfig } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let teams: data.Team[];

  // Note: https://iro-color.com/colorchart/tone/bright-tone.html
  const CHART_COLORS = ["#00A95F", "#EA5532"];
  const CHART_INFO: { label: string; func: GetStatsFunction }[] = [
    {
      label: "PR(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.ship.pr;
      },
    },
    {
      label: "勝率(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.ship.win_rate;
      },
    },
    {
      label: "ダメージ(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.ship.damage;
      },
    },
    {
      label: "戦闘数(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.ship.battles;
      },
    },
    {
      label: "PR(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.pr;
      },
    },
    {
      label: "勝率(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.win_rate;
      },
    },
    {
      label: "ダメージ(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.damage;
      },
    },
    {
      label: "戦闘数(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.battles;
      },
    },
  ];

  const CHART_INFO_OVERALL: { label: string; func: GetStatsFunction }[] = [
    {
      label: "PR(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.pr;
      },
    },
    {
      label: "勝率(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.win_rate;
      },
    },
    {
      label: "ダメージ(総合)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.damage;
      },
    },
    {
      label: "戦闘数(総合)",
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

<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
  {#each CHART_INFO as item}
    {@const max = getMaxValueInAllPlayers(teams, item.func)}

    <table
      class="charts-css column multiple show-labels data-spacing-4 datasets-spacing-4 show-heading"
    >
      <tbody class="h-24">
        <tr>
          <th scope="row" class="text-sm">{item.label}</th>
          {#each teams as team, i}
            {@const values = getValuesInTeam(team, item.func)}
            {@const mn = calculateMean(values)}
            {@const sd = calculateStandardDeviation(values)}
            {@const isTranslateText = mn / max <= 0.5}
            {@const textColorClass = isTranslateText ? "" : "text-neutral-100"}

            <td style="--size: calc({mn}/{max}); --color: {CHART_COLORS[i]}">
              <div
                class="py-1 text-center leading-none {isTranslateText
                  ? 'absolute -translate-y-full'
                  : ''}"
              >
                <div
                  class="leading-none text-sm text-nowrap {textColorClass} font-semibold"
                >
                  {formatWithSuffix(mn)}
                </div>
                <div class="leading-none text-sm text-nowrap {textColorClass}">
                  (±{formatWithSuffix(sd)})
                </div>
              </div>
            </td>
          {/each}
        </tr>
      </tbody>
    </table>
  {/each}
</div>
