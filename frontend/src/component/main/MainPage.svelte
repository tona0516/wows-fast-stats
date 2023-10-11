<script lang="ts">
  import { Screenshot } from "src/lib/Screenshot";
  import { MAIN_PAGE_ID, ScreenshotType } from "src/lib/types";
  import BattleMeta from "src/component/main/internal/BattleMeta.svelte";
  import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
  import {
    storedBattle,
    storedUserConfig,
    storedSummary,
    storedRequiredConfigError,
  } from "src/stores";
  import Function from "./internal/Function.svelte";
  import Summary from "./internal/Summary.svelte";
  import ColorDescription from "./internal/ColorDescription.svelte";
  import Ofuse from "./internal/Ofuse.svelte";
  import UkSpinner from "../common/uikit/UkSpinner.svelte";
  import { Battle } from "wailsjs/go/main/App";
  import { createEventDispatcher } from "svelte";

  const screenshot = new Screenshot();
  const dispatch = createEventDispatcher();

  let isLoading = false;

  export const fetchBattle = async () => {
    try {
      isLoading = true;

      // Note: 過去のデータが影響してか値が0になってしまうためクリーンする
      storedBattle.set(undefined);

      const start = new Date().getTime();
      const battle = await Battle();
      const elapsed = (new Date().getTime() - start) / 1000;

      storedBattle.set(battle);
      dispatch("FetchSuccess", {
        message: `データ取得完了: ${elapsed.toFixed(1)}秒`,
      });

      if ($storedUserConfig.save_screenshot) {
        screenshot.take(ScreenshotType.auto, battle.meta);
      }
    } catch (error) {
      dispatch("Failure", error);
    } finally {
      isLoading = false;
    }
  };
</script>

<!-- Note: Use the same color as that of body.  -->
<div
  id={MAIN_PAGE_ID}
  class="uk-padding-small uk-light uk-background-secondary"
>
  <div class="uk-margin-small uk-flex uk-flex-center">
    <Function {screenshot} on:ScreenshotSaved on:Failure />
  </div>

  <div class="uk-margin-small">
    {#if $storedBattle}
      {@const teams = $storedBattle.teams}
      {@const meta = $storedBattle.meta}
      {@const userConfig = $storedUserConfig}

      <div class="uk-flex uk-flex-center">
        <StatisticsTable
          {teams}
          {userConfig}
          on:EditAlertPlayer
          on:RemoveAlertPlayer
          on:CheckPlayer
        />
      </div>

      <div class="uk-flex uk-flex-center">
        <BattleMeta {meta} />
      </div>

      {#if $storedSummary}
        {@const summary = $storedSummary}
        <div class="uk-flex uk-flex-center">
          <Summary {summary} />
        </div>
      {/if}

      <ColorDescription config={userConfig} />
    {:else}
      <p class="uk-text-center">
        {#if $storedRequiredConfigError.valid}
          戦闘中ではありません。開始時に自動的にリロードします。
        {:else}
          設定画面から初期設定を行ってください。
        {/if}
      </p>
    {/if}
  </div>

  <div class="uk-margin-small">
    <Ofuse />
  </div>

  {#if isLoading}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <UkSpinner />
      </div>
    </div>
  {/if}
</div>
