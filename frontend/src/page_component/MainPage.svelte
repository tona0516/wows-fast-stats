<script lang="ts">
  import { format, fromUnixTime } from "date-fns";
  import { MAIN_PAGE_ID } from "src/const";
  import ColorDescription from "src/other_component/ColorDescription.svelte";
  import Ofuse from "src/other_component/Ofuse.svelte";
  import StatisticsTable from "src/other_component/StatisticsTable.svelte";
  import { storedBattle, storedUserConfig, storedSummary } from "src/stores";

  const formattedDate = (unixtime: number): string => {
    return format(fromUnixTime(unixtime), "yyyy/MM/dd HH:mm:ss");
  };
</script>

<div id={MAIN_PAGE_ID}>
  {#if $storedBattle}
    <div class="mt-2 mx-2">
      <StatisticsTable
        teams={$storedBattle.teams}
        userConfig={$storedUserConfig}
        on:UpdateAlertPlayer
        on:RemoveAlertPlayer
        on:CheckPlayer
      />
    </div>

    {#if $storedSummary}
      <div class="d-flex justify-content-center">
        <div class="center mx-2">
          <table class="table table-sm table-text-color w-auto td-multiple">
            <tbody>
              <tr>
                <td>日時</td>
                <td>{formattedDate($storedBattle.meta.unixtime)}</td>
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
                <th colspan={$storedSummary.shipColspan}>艦成績</th>
                <th colspan={$storedSummary.overallColspan}>総合成績</th>
              </tr>
              <tr>
                <th />
                {#each $storedSummary.labels as label}
                  <th>{label}</th>
                {/each}
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>{$storedBattle.teams[0].name}</td>
                {#each $storedSummary.friends as friend}
                  <td>{friend}</td>
                {/each}
              </tr>

              <tr>
                <td>{$storedBattle.teams[1].name}</td>
                {#each $storedSummary.enemies as enemy}
                  <td>{enemy}</td>
                {/each}
              </tr>

              <tr>
                <td>差</td>
                {#each $storedSummary.diffs as diff}
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
</div>

<Ofuse />
