<script lang="ts">
  import type { Screenshot } from "src/lib/Screenshot";
  import { MAIN_PAGE_ID } from "src/lib/types";
  import BattleMeta from "src/other_component/BattleMeta.svelte";
  import ColorDescription from "src/other_component/ColorDescription.svelte";
  import Function from "src/other_component/Function.svelte";
  import Ofuse from "src/other_component/Ofuse.svelte";
  import StatisticsTable from "src/other_component/StatisticsTable.svelte";
  import Summary from "src/other_component/Summary.svelte";
  import { storedBattle, storedUserConfig, storedSummary } from "src/stores";

  export let isLoading: boolean;
  export let screenshot: Screenshot;
</script>

<!-- Note: Use the same color as that of body.  -->
<div
  id={MAIN_PAGE_ID}
  class="uk-padding-small uk-light uk-background-secondary"
>
  <div class="uk-margin-small">
    <Function {screenshot} on:ScreenshotSaved on:Failure />
  </div>

  <div class="uk-margin-small">
    {#if $storedBattle}
      <StatisticsTable
        teams={$storedBattle.teams}
        userConfig={$storedUserConfig}
        on:EditAlertPlayer
        on:RemoveAlertPlayer
        on:CheckPlayer
      />

      <BattleMeta meta={$storedBattle.meta} />

      {#if $storedSummary}
        <Summary summary={$storedSummary} />
      {/if}

      <ColorDescription userConfig={$storedUserConfig} />
    {:else}
      <div class="uk-text-center">
        <p>戦闘中ではありません。開始時に自動的にリロードします。</p>
      </div>
    {/if}
  </div>

  <div class="uk-margin-small">
    <Ofuse />
  </div>

  {#if isLoading}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <div uk-spinner />
      </div>
    </div>
  {/if}
</div>
