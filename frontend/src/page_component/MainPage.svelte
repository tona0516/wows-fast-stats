<script lang="ts">
import { get } from "svelte/store";
import {
  storedBattle,
  storedUserConfig,
  storedSummaryResult,
  storedAlertPlayers,
} from "../stores";
import ColorDescription from "../other_component/ColorDescription.svelte";
import Ofuse from "../other_component/Ofuse.svelte";
import StatisticsTable from "../other_component/StatisticsTable.svelte";

let battle = get(storedBattle);
storedBattle.subscribe((it) => (battle = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));

let alertPlayers = get(storedAlertPlayers);
storedAlertPlayers.subscribe((it) => (alertPlayers = it));

let summaryResult = get(storedSummaryResult);
storedSummaryResult.subscribe((it) => {
  summaryResult = it;
});
</script>

{#if battle}
  <div class="m-2">
    <StatisticsTable
      teams="{battle.teams}"
      userConfig="{userConfig}"
      alertPlayers="{alertPlayers}"
      on:UpdateAlertPlayer
      on:RemoveAlertPlayer
      on:CheckPlayer
    />
  </div>

  {#if summaryResult}
    <div class="d-flex justify-content-center">
      <div class="center mx-2">
        <h6>戦闘情報</h6>

        <table class="table table-sm table-text-color w-auto td-multiple">
          <tbody>
            <tr>
              <td>日時</td>
              <td>{battle.meta.date}</td>
            </tr>

            <tr>
              <td>戦闘タイプ</td>
              <td>{battle.meta.type}</td>
            </tr>

            <tr>
              <td>マップ</td>
              <td>{battle.meta.arena}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="center mx-2">
        <h6>チーム平均</h6>

        <table class="table table-sm table-text-color w-auto td-multiple">
          <thead>
            <tr>
              <th></th>
              <th colspan="{summaryResult.shipColspan}">艦成績</th>
              <th colspan="{summaryResult.overallColspan}">総合成績</th>
            </tr>
            <tr>
              <th></th>
              {#each summaryResult.labels as label}
                <th>{label}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{battle.teams[0].name}</td>
              {#each summaryResult.friends as friend}
                <td>{friend}</td>
              {/each}
            </tr>

            <tr>
              <td>{battle.teams[1].name}</td>
              {#each summaryResult.enemies as enemy}
                <td>{enemy}</td>
              {/each}
            </tr>

            <tr>
              <td>差</td>
              {#each summaryResult.diffs as diff}
                <td class="{diff.colorClass}">{diff.value}</td>
              {/each}
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  {/if}
{/if}

<ColorDescription />

<Ofuse />
