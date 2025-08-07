<script lang="ts">
  import { max } from "date-fns";
  import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
  import type { StatsExtra } from "src/lib/types";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;
  export let statsExtra: StatsExtra;

  $: teamThreatLevels = TeamThreatLevel.fromBattle(battle, statsExtra);
</script>

{#if teamThreatLevels}
  {@const friendRate =
    teamThreatLevels[0].average /
    teamThreatLevels.map((l) => l.average).reduce((a, b) => a + b)}
  <div>
    <table class="charts-css bar stacked">
      <tbody>
        <tr>
          <td style="--size: calc({friendRate}); --color: #00A95F;"></td>
          <td style="--size: calc{1 - friendRate}); --color: #EA5532;"></td>
        </tr>
      </tbody>
    </table>

    {#each teamThreatLevels as level}
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-title">戦力評価</div>
          <div class="stat-value">{level.average.format(0)}</div>
        </div>

        <div class="stat">
          <div class="stat-title">精度</div>
          <div class="stat-value">{level.accuracy.format(0)}%</div>
        </div>

        <div class="stat">
          <div class="stat-title">介護指数</div>
          <div class="stat-value">{level.dissociationDegree.format(0)}%</div>
        </div>
      </div>
    {/each}
  </div>
{/if}
