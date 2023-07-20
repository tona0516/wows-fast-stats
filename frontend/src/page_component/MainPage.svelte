<script lang="ts">
  import {
    storedBattle,
    storedUserConfig,
    storedSummaryResult,
    storedAlertPlayers,
    storedLogs,
  } from "../stores";
  import ColorDescription from "../other_component/ColorDescription.svelte";
  import Ofuse from "../other_component/Ofuse.svelte";
  import StatisticsTable from "../other_component/StatisticsTable.svelte";
  import { toDateForDisplay } from "../util";
  import Logging from "../other_component/Logging.svelte";
</script>

{#if $storedBattle}
  <div class="mt-2 mx-2">
    <StatisticsTable
      teams={$storedBattle.teams}
      userConfig={$storedUserConfig}
      alertPlayers={$storedAlertPlayers}
      on:UpdateAlertPlayer
      on:RemoveAlertPlayer
      on:CheckPlayer
    />
  </div>

  {#if $storedSummaryResult}
    <div class="d-flex justify-content-center">
      <div class="center mx-2">
        <table class="table table-sm table-text-color w-auto td-multiple">
          <tbody>
            <tr>
              <td>日時</td>
              <td>{toDateForDisplay($storedBattle.meta.unixtime)}</td>
            </tr>

            <tr>
              <td>戦闘タイプ</td>
              <td>{$storedBattle.meta.type}</td>
            </tr>

            <tr>
              <td>マップ</td>
              <td>{$storedBattle.meta.arena}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="center mx-2">
        <table class="table table-sm table-text-color w-auto td-multiple">
          <thead>
            <tr>
              <th />
              <th colspan={$storedSummaryResult.shipColspan}>艦成績</th>
              <th colspan={$storedSummaryResult.overallColspan}>総合成績</th>
            </tr>
            <tr>
              <th />
              {#each $storedSummaryResult.labels as label}
                <th>{label}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{$storedBattle.teams[0].name}</td>
              {#each $storedSummaryResult.friends as friend}
                <td>{friend}</td>
              {/each}
            </tr>

            <tr>
              <td>{$storedBattle.teams[1].name}</td>
              {#each $storedSummaryResult.enemies as enemy}
                <td>{enemy}</td>
              {/each}
            </tr>

            <tr>
              <td>差</td>
              {#each $storedSummaryResult.diffs as diff}
                <td class={diff.colorClass}>{diff.value}</td>
              {/each}
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  {/if}
{/if}

<ColorDescription userConfig={$storedUserConfig} />

<Logging logs={$storedLogs} />

<Ofuse />
