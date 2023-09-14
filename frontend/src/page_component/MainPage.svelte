<script lang="ts">
  import { format, fromUnixTime } from "date-fns";
  import { MAIN_PAGE_ID } from "src/const";
  import ColorDescription from "src/other_component/ColorDescription.svelte";
  import Ofuse from "src/other_component/Ofuse.svelte";
  import StatisticsTable from "src/other_component/StatisticsTable.svelte";
  import SummaryGraph from "src/other_component/SummaryGraph.svelte";
  import { storedBattle, storedUserConfig, storedSummary } from "src/stores";

  const formattedDate = (unixtime: number): string => {
    return format(fromUnixTime(unixtime), "yyyy/MM/dd HH:mm:ss");
  };

  let battleMetas: {
    icon: string;
    text: string;
  }[] = [];
  $: {
    if ($storedBattle) {
      battleMetas = [
        {
          icon: "bi bi-clock-fill",
          text: formattedDate($storedBattle.meta.unixtime),
        },
        { icon: "bi bi-tag-fill", text: $storedBattle.meta.type },
        { icon: "bi bi-geo-alt-fill", text: $storedBattle.meta.arena },
      ];
    } else {
      battleMetas = [];
    }
  }
</script>

<div id={MAIN_PAGE_ID}>
  {#if battleMetas}
    <div class="d-flex mt-1 mx-2">
      {#each battleMetas as meta}
        <div class="me-2">
          <i class={meta.icon}></i>{meta.text}
        </div>
      {/each}
    </div>
  {/if}

  {#if $storedBattle}
    <div class="mt-1 mx-2 center">
      <StatisticsTable
        teams={$storedBattle.teams}
        userConfig={$storedUserConfig}
        on:UpdateAlertPlayer
        on:RemoveAlertPlayer
        on:CheckPlayer
      />
    </div>
  {/if}

  {#if $storedSummary}
    <div class="mt-1 mx-2">
      <SummaryGraph summary={$storedSummary} />
    </div>
  {/if}

  <ColorDescription userConfig={$storedUserConfig} />
</div>

<Ofuse />
