<script lang="ts">
  import { TeamThreatLevel } from "src/lib/TeamThreatLevel";
  import type { StatsExtra } from "src/lib/types";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;
  export let statsExtra: StatsExtra;

  const teamThreatLevels = TeamThreatLevel.fromBattle(battle, statsExtra);
</script>

<div>
  {#each teamThreatLevels as level}
  <div class="flex flex-col items-center">{level[0]}</div>
    <div class="stats shadow">
      <div class="stat">
        <div class="stat-title">戦力評価</div>
        <div class="stat-value">{level[1].average.format(0)}</div>
      </div>

      <div class="stat">
        <div class="stat-title">精度</div>
        <div class="stat-value">{level[1].accuracy.format(0)}%</div>
      </div>

      <div class="stat">
        <div class="stat-title">介護指数</div>
        <div class="stat-value">{level[1].dissociationDegree.format(0)}%</div>
      </div>
    </div>
  {/each}
</div>
