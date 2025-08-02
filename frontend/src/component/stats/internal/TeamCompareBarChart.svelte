<script lang="ts">
  import { mean, standardDeviation } from "simple-statistics";
  import type { GetStatsFunction } from "src/lib/types";
  import { formatWithSuffix, toPlayerStats } from "src/lib/util";
  import { storedConfig } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;

  export let caption = "全艦艇";
  export let filterFunc = (player: data.Player) => {
    return true;
  };

  $: teams = filterPlayers(battle.teams);

  // Note: https://iro-color.com/colorchart/tone/bright-tone.html
  const CHART_COLORS = ["#00A95F", "#EA5532"];
  const CHART_INFO: { label: string; func: GetStatsFunction }[] = [
    {
      label: "PR(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.pr;
      },
    },
    {
      label: "ダメージ(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.damage;
      },
    },
    {
      label: "勝率(艦)",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.win_rate;
      },
    },
    {
      label: "戦力評価",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.threat_level.modified;
      },
    },
    {
      label: "戦闘数",
      func: (ps: data.PlayerStats): number => {
        return ps.overall.battles;
      },
    },
  ];

  function filterPlayers(team: data.Team[]): data.Team[] {
    return team.map((team) => {
      const updatedTeam = structuredClone(team);
      updatedTeam.players = team.players.filter(filterFunc);
      return updatedTeam;
    });
  }

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

<div class="w-2xl">
  <table
    class="charts-css column multiple show-labels data-spacing-4 datasets-spacing-4 show-heading"
  >
    <caption>{caption}</caption>
    {#if teams.every((team) => team.players.length > 0)}
      <tbody class="h-24">
        {#each CHART_INFO as item}
          {@const max = getMaxValueInAllPlayers(teams, item.func)}
          <tr>
            <th scope="row" class="text-sm">{item.label}</th>
            {#each teams as team, i}
              {@const values = getValuesInTeam(team, item.func)}
              {@const mn = calculateMean(values)}
              {@const sd = calculateStandardDeviation(values)}
              {@const isTranslateText = mn / max <= 0.5}
              {@const textColorClass = isTranslateText
                ? ""
                : "text-neutral-100"}

              <td style="--size: calc({mn}/{max}); --color: {CHART_COLORS[i]}">
                <div
                  class="py-1 text-center leading-none {isTranslateText
                    ? 'absolute -translate-y-full'
                    : ''}"
                >
                  <span
                    class="leading-none text-sm text-nowrap {textColorClass} font-semibold"
                  >
                    {formatWithSuffix(mn)}
                  </span>
                  <span
                    class="leading-none text-sm text-nowrap {textColorClass}"
                  >
                    (±{formatWithSuffix(sd)})
                  </span>
                </div>
              </td>
            {/each}
          </tr>
        {/each}
      </tbody>
    {:else}
      <div class="inline-block align-middle">
        <div role="alert" class="alert alert-warning alert-dash flex flex-col">
          <span class="items-center">対象艦種なし</span>
        </div>
      </div>
    {/if}
  </table>
</div>
